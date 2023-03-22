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
	return ro.fs.Open(ConvertAbs(name))
}

func (ro rofs) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	if flag != os.O_RDONLY {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
	}

	f, err := ro.Open(name)
	if err != nil {
		return nil, err
	}

	return rofile{f}, nil
}

func (ro rofs) Remove(name string) error {
	return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrPermission}
}

func (ro rofs) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(ro.fs, ConvertAbs(name))
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

type rofile struct{ fs.File }

var _ File = rofile{}

// WrapFile wraps a read-only file into a read-writable file. Any functions that
// write to the file will return an error.
func WrapFile(f fs.File) File {
	return rofile{f}
}

func (f rofile) Write([]byte) (int, error) {
	return 0, fs.ErrPermission
}
