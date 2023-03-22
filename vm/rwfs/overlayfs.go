package rwfs

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

// OverlayFS is a read-writable filesystem that overlays multiple filesystems.
func OverlayFS(rw FS, ro fs.FS) FS {
	return overlayFS{rw, ro}
}

type overlayFS struct {
	rw FS
	ro fs.FS
}

var (
	_ FS           = overlayFS{}
	_ fs.ReadDirFS = overlayFS{}
)

func (o overlayFS) Open(name string) (fs.File, error) {
	name = ConvertAbs(name)

	// Check the read-only filesystem first because it's read-only, so it cannot
	// be removed or written to.
	if f, err := o.ro.Open(name); err == nil {
		return readDirableROFile{f, name, o}, nil
	}

	// Then check the read-write filesystem.
	if f, err := o.rw.Open(name); err == nil {
		return readDirableROFile{f, name, o}, nil
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

func (o overlayFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	name = ConvertAbs(name)

	f, err := o.ro.Open(name)
	if err == nil {
		if flag&os.O_RDONLY != 0 {
			return rofile{f}, nil
		}
		return nil, &fs.PathError{
			Op:   "openfile",
			Path: name,
			Err:  fs.ErrPermission,
		}
	}

	// If we're creating a file, then there's a chance that the directory that
	// we're in exists on the read-only filesystem but not the read-write
	// filesystem. In that case, we'll need to create the directory.
	if flag&os.O_CREATE != 0 {
		dir := path.Dir(name)

		s, err := fs.Stat(o.ro, dir)
		if err == nil {
			if err := o.rw.MkdirAll(dir, s.Mode()); err != nil {
				return nil, &fs.PathError{
					Op:   "openfile",
					Path: name,
					Err:  err,
				}
			}
		}
	}

	return o.rw.OpenFile(name, flag, perm)
}

func (o overlayFS) Remove(name string) error {
	name = ConvertAbs(name)

	err := o.rw.Remove(name)
	// Either no error or does-not-exist error are fine. It's the two common and
	// valid errors that we can assume RW can return.
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	// Consider every other error a permission error, since we'd have to check
	// the read-only filesystem otherwise.
	return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrPermission}
}

func (o overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = ConvertAbs(name)

	// For ReadDir, combine.
	rwEntries, err1 := fs.ReadDir(o.rw, name)
	roEntries, err2 := fs.ReadDir(o.ro, name)

	// If the directory exists in the read-only filesystem, then we're fine with
	// the read-write filesystem returning an error.
	if err2 == nil && err1 != nil {
		return roEntries, nil
	}

	entries := append(rwEntries, roEntries...)
	// Deduplicate paths, because we're handling writes to a read-only directory
	// by making it on the read-write filesystem as well.
	entries = DeduplicateDirEntries(entries)
	return entries, errors.Join(err1, err2)
}

func (o overlayFS) Mkdir(name string, perm fs.FileMode) error {
	name = ConvertAbs(name)

	// Check if the directory already exists on either filesystems.
	_, err1 := fs.Stat(o.rw, name)
	_, err2 := fs.Stat(o.ro, name)
	if err1 == nil || err2 == nil {
		return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrExist}
	}

	// We can only do this on a read-write filesystem.
	return o.rw.Mkdir(name, perm)
}

func (o overlayFS) MkdirAll(name string, perm fs.FileMode) error {
	name = ConvertAbs(name)

	// Same as above.
	_, err1 := fs.Stat(o.rw, name)
	_, err2 := fs.Stat(o.ro, name)
	if err1 == nil || err2 == nil {
		return nil
	}

	return o.rw.MkdirAll(name, perm)
}

func (o overlayFS) RemoveAll(name string) error {
	name = ConvertAbs(name)

	// Does the directory exist on the read-only filesystems?
	_, err1 := fs.Stat(o.rw, name)
	if err1 == nil {
		return &fs.PathError{Op: "removeall", Path: name, Err: fs.ErrPermission}
	}

	// Otherwise, the error is whatever the read-write filesystem returns.
	return o.rw.RemoveAll(name)
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
