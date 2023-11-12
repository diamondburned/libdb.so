package rwfs

import (
	"io/fs"
	"os"
)

// ReadOnlyFS wraps a read-only filesystem into a read-writable filesystem. Any
// functions that write to the filesystem will return an error.
func ReadOnlyFS(fs fs.FS) FS {
	return rofs{fs}
}

type rofs struct{ fs fs.FS }

var _ FS = rofs{}

func (ro rofs) Open(name string) (fs.File, error) {
	f, err := ro.fs.Open(ConvertAbs(name))
	if err != nil {
		return nil, err
	}

	return roFile{f}, nil
}

func (ro rofs) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	if flag != os.O_RDONLY {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
	}

	f, err := ro.Open(ConvertAbs(name))
	if err != nil {
		return nil, err
	}

	return roFile{f}, nil
}

func (ro rofs) Remove(name string) error {
	return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrPermission}
}

func (ro rofs) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(ro.fs, ConvertAbs(name))
	if err != nil {
		return nil, err
	}

	for i := range entries {
		entries[i] = roDirEntry{entries[i]}
	}

	return entries, nil
}

func (ro rofs) Mkdir(name string, perm fs.FileMode) error {
	return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrPermission}
}

func (ro rofs) MkdirAll(name string, perm fs.FileMode) error {
	return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrPermission}
}

func (ro rofs) RemoveAll(name string) error {
	return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrPermission}
}

type roFile struct{ fs.File }

var _ File = roFile{}

// WrapFile wraps a read-only file into a read-writable file. Any functions that
// write to the file will return an error.
func WrapFile(f fs.File) File {
	return roFile{f}
}

func (f roFile) Write([]byte) (int, error) {
	return 0, fs.ErrPermission
}

func (f roFile) Stat() (fs.FileInfo, error) {
	s, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return roStat{s}, nil
}

type roStat struct{ fs.FileInfo }

var _ fs.FileInfo = roStat{}

func (s roStat) Mode() fs.FileMode {
	return s.FileInfo.Mode() &^ 0222
}

type roDirEntry struct{ fs.DirEntry }

var _ fs.DirEntry = roDirEntry{}

func (e roDirEntry) Type() fs.FileMode {
	return e.DirEntry.Type() &^ 0222
}

func (e roDirEntry) Info() (fs.FileInfo, error) {
	s, err := e.DirEntry.Info()
	if err != nil {
		return nil, err
	}
	return roStat{s}, nil
}
