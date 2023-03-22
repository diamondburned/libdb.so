package rwfs

import (
	"errors"
	"io/fs"

	"libdb.so/vm/rwfs/internal/fsutil"
)

// OverlayFS is a read-writable filesystem that overlays multiple filesystems.
func OverlayFS(mounts ...FS) FS {
	return overlayFS(mounts)
}

type overlayFS []FS

var (
	_ FS           = overlayFS{}
	_ fs.ReadDirFS = overlayFS{}
)

func (ofs overlayFS) Open(name string) (fs.File, error) {
	name = fsutil.ConvertAbs(name)

	var errs []error

	for _, vfs := range ofs {
		f, err := vfs.Open(name)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			errs = append(errs, err)
			continue
		}
		return readDirableROFile{f, name, ofs}, nil
	}

	return nil, errorsOrNotFound(errs)
}

func (ofs overlayFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	name = fsutil.ConvertAbs(name)

	var errs []error

	for _, vfs := range ofs {
		f, err := vfs.OpenFile(name, flag, perm)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			errs = append(errs, err)
			continue
		}
		return readDirableRWFile{f, name, ofs}, nil
	}

	return nil, errorsOrNotFound(errs)
}

func (ofs overlayFS) Remove(name string) error {
	name = fsutil.ConvertAbs(name)

	var errs []error

	for _, vfs := range ofs {
		if err := vfs.Remove(name); err != nil && !errors.Is(err, fs.ErrNotExist) {
			errs = append(errs, err)
		}
	}

	return errorsOrNotFound(errs)
}

func (ofs overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = fsutil.ConvertAbs(name)

	var allEntries []fs.DirEntry
	var errs []error

	for _, vfs := range ofs {
		entries, err := fs.ReadDir(vfs, name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		allEntries = append(allEntries, entries...)
	}

	return allEntries, errors.Join(errs...)
}

func (ofs overlayFS) Mkdir(name string, perm fs.FileMode) error {
	name = fsutil.ConvertAbs(name)

	var errs []error

	for _, vfs := range ofs {
		err := vfs.Mkdir(name, perm)
		if err != nil {
			if errors.Is(err, fs.ErrExist) {
				continue
			}
			errs = append(errs, err)
			continue
		}
		return nil
	}

	return errors.Join(errs...)
}

func (ofs overlayFS) MkdirAll(name string, perm fs.FileMode) error {
	name = fsutil.ConvertAbs(name)

	var errs []error

	for _, vfs := range ofs {
		if err := vfs.MkdirAll(name, perm); err != nil {
			errs = append(errs, err)
			continue
		}
		return nil
	}

	return errors.Join(errs...)
}

func (ofs overlayFS) RemoveAll(name string) error {
	name = fsutil.ConvertAbs(name)

	var errs []error

	for _, vfs := range ofs {
		if err := vfs.RemoveAll(name); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
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
