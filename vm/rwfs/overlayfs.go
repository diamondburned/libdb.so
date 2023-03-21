package rwfs

import (
	"errors"
	"io/fs"
)

// OverlayFS is a read-writable filesystem that overlays multiple filesystems.
func OverlayFS(mounts ...FS) FS {
	return overlayFS(mounts)
}

type overlayFS []FS

var _ FS = overlayFS{}

func (ofs overlayFS) Open(name string) (fs.File, error) {
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
		return f, nil
	}

	return nil, errors.Join(errs...)
}

func (ofs overlayFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
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
		return f, nil
	}

	return nil, errors.Join(errs...)
}

func (ofs overlayFS) Remove(name string) error {
	var errs []error

	for _, vfs := range ofs {
		if err := vfs.Remove(name); err != nil && !errors.Is(err, fs.ErrNotExist) {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (ofs overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	var errs []error

	for _, vfs := range ofs {
		entries, err := fs.ReadDir(vfs, name)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			errs = append(errs, err)
			continue
		}
		entries = append(entries, entries...)
	}

	return entries, errors.Join(errs...)
}

func (ofs overlayFS) Mkdir(name string, perm fs.FileMode) error {
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
	var errs []error

	for _, vfs := range ofs {
		if err := vfs.RemoveAll(name); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
