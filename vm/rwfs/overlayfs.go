package rwfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

// OverlayFS is a read-writable filesystem that overlays multiple filesystems.
func OverlayFS(rw FS, ro ...fs.FS) FS {
	return overlayFS{rw, ro}
}

type overlayFS struct {
	rw FS
	ro []fs.FS
}

var (
	_ FS           = overlayFS{}
	_ fs.ReadDirFS = overlayFS{}
)

func (o overlayFS) Open(name string) (fs.File, error) {
	name = ConvertAbs(name)

	// Check the read-only filesystem first because it's read-only, so it cannot
	// be removed or written to.
	for _, ro := range o.ro {
		if f, err := ro.Open(name); !errors.Is(err, fs.ErrNotExist) {
			if err != nil {
				return nil, err
			}
			return readDirableROFile{f, name, o}, nil
		}
	}

	// Then check the read-write filesystem.
	if f, err := o.rw.Open(name); err == nil {
		return readDirableROFile{f, name, o}, nil
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

func (o overlayFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	name = ConvertAbs(name)

	for _, ro := range o.ro {
		f, err := ro.Open(name)
		if err == nil {
			if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND) != 0 {
				f.Close()
				return nil, &fs.PathError{
					Op:   "openfile",
					Path: name,
					Err:  fs.ErrPermission,
				}
			}
			return roFile{f}, nil
		} else {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
		}
	}

	// If we're creating a file, then there's a chance that the directory that
	// we're in exists on the read-only filesystem but not the read-write
	// filesystem. In that case, we'll need to create the directory.
	if flag&os.O_CREATE != 0 {
		dir := path.Dir(name)

		for _, ro := range o.ro {
			s, err := fs.Stat(ro, dir)
			if err == nil {
				if !s.IsDir() {
					return nil, &fs.PathError{
						Op:   "openfile",
						Path: name,
						Err:  fs.ErrInvalid,
					}
				}

				if err := o.rw.MkdirAll(dir, s.Mode()); err != nil {
					return nil, &fs.PathError{
						Op:   "openfile",
						Path: name,
						Err:  err,
					}
				}

				break
			} else {
				if !errors.Is(err, fs.ErrNotExist) {
					return nil, err
				}
			}
		}
	}

	return o.rw.OpenFile(name, flag, perm)
}

func (o overlayFS) Remove(name string) error {
	name = ConvertAbs(name)

	if err := o.rw.Remove(name); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		// Don't allow if we're removing a read-only file or directory.
		for _, ro := range o.ro {
			_, err1 := fs.Stat(ro, name)
			if err1 == nil {
				return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrPermission}
			}
		}
	}

	return nil
}

func (o overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = ConvertAbs(name)

	var entries []fs.DirEntry
	for i, ro := range o.ro {
		roEntries, err := fs.ReadDir(ro, name)
		if err == nil {
			entries = append(entries, roEntries...)
		} else {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("error at fs %d: %w", i, err)
			}
		}
	}

	rwEntries, err := fs.ReadDir(o.rw, name)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("error at rw fs: %w", err)
	}
	entries = append(entries, rwEntries...)

	// Deduplicate paths, because we're handling writes to a read-only directory
	// by making it on the read-write filesystem as well.
	entries = DeduplicateDirEntries(entries)
	return entries, nil
}

func (o overlayFS) Mkdir(name string, perm fs.FileMode) error {
	name = ConvertAbs(name)

	// Check if the directory already exists on either filesystems.
	_, err := fs.Stat(o.rw, name)
	if errors.Is(err, fs.ErrExist) {
		return err
	}
	for _, ro := range o.ro {
		_, err := fs.Stat(ro, name)
		if errors.Is(err, fs.ErrExist) {
			return err
		}
	}

	// We can only do this on a read-write filesystem.
	return o.rw.Mkdir(name, perm)
}

func (o overlayFS) MkdirAll(name string, perm fs.FileMode) error {
	name = ConvertAbs(name)

	return o.rw.MkdirAll(name, perm)
}

func (o overlayFS) RemoveAll(name string) error {
	name = ConvertAbs(name)

	// Does the directory exist on the read-only filesystems?
	for _, ro := range o.ro {
		_, err1 := fs.Stat(ro, name)
		if err1 == nil {
			return &fs.PathError{Op: "removeall", Path: name, Err: fs.ErrPermission}
		}
	}

	// Otherwise, the error is whatever the read-write filesystem returns.
	return o.rw.RemoveAll(name)

	// Note to self: I think rm -f will actually try to remove RW files off a RO
	// directory. Maybe rm itself implements this logic?
}

func errorsOrNotFound(errs []error) error {
	if len(errs) == 0 {
		return fs.ErrNotExist
	}
	return errors.Join(errs...)
}

var _ fs.ReadDirFile = readDirableROFile{}
var _ fs.ReadDirFile = readDirableRWFile{}

type readDirableROFile struct {
	fs.File
	path    string
	overlay overlayFS
}

func (f readDirableROFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if _, ok := f.File.(fs.ReadDirFile); !ok {
		return nil, fs.ErrInvalid
	}
	return f.overlay.ReadDir(f.path)
}

type readDirableRWFile struct {
	File
	path    string
	overlay overlayFS
}

func (f readDirableRWFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if _, ok := f.File.(fs.ReadDirFile); !ok {
		return nil, fs.ErrInvalid
	}
	return f.overlay.ReadDir(f.path)
}
