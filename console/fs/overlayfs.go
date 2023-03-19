package fs

import (
	"errors"
	"io/fs"
)

// OverlayFS is a read-writable filesystem that overlays multiple filesystems.
func OverlayFS(mounts ...FS) FS {
	return overlayFS(mounts)
}

type overlayFS []FS

var (
	_ FS    = overlayFS{}
	_ DirFS = overlayFS{}
)

func (ofs overlayFS) Open(name string) (fs.File, error) {
	for _, vfs := range ofs {
		f, err := vfs.Open(name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		return f, err
	}

	return nil, fs.ErrNotExist
}

func (ofs overlayFS) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	for _, vfs := range ofs {
		f, err := vfs.OpenFile(name, flag, perm)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		return f, err
	}

	return nil, fs.ErrNotExist
}

func (ofs overlayFS) Remove(name string) error {
	for _, vfs := range ofs {
		_, err := fs.Stat(vfs, name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		return vfs.Remove(name)
	}

	return nil
}

func (ofs overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	for _, vfs := range ofs {
		entries, err := fs.ReadDir(vfs, name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}

		return entries, err
	}

	return nil, fs.ErrNotExist
}

func (ofs overlayFS) Mkdir(name string, perm fs.FileMode) error {
	for _, vfs := range ofs {
		dirfs, ok := vfs.(DirFS)
		if !ok {
			continue
		}

		err := dirfs.Mkdir(name, perm)
		if errors.Is(err, fs.ErrExist) {
			continue
		}

		return err
	}

	return nil
}

func (ofs overlayFS) MkdirAll(name string, perm fs.FileMode) error {
	for _, vfs := range ofs {
		dirfs, ok := vfs.(DirFS)
		if !ok {
			continue
		}

		err := dirfs.MkdirAll(name, perm)
		return err
	}

	return nil
}

func (ofs overlayFS) RemoveAll(name string) error {
	for _, vfs := range ofs {
		dirfs, ok := vfs.(DirFS)
		if !ok {
			continue
		}

		err := dirfs.RemoveAll(name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}

		return err
	}

	return nil
}
