package fs

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

// ErrReadOnly is returned when a filesystem is read-only.
var ErrReadOnly = errors.New("read-only filesystem")

// ReadOnlyFS wraps a read-only filesystem into a read-writable filesystem. Any
// functions that write to the filesystem will return an error.
func ReadOnlyFS(fs fs.FS) FS {
	return rofs{fs}
}

func fromAbs(abs string) string {
	name := abs
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		name = "."
	}
	return name
}

type rofs struct{ fs fs.FS }

var (
	_ FS    = rofs{}
	_ DirFS = rofs{}
)

func (ro rofs) Open(name string) (fs.File, error) {
	return ro.fs.Open(fromAbs(name))
}

func (ro rofs) OpenFile(name string, flag int, perm fs.FileMode) (File, error) {
	if flag != os.O_RDONLY {
		return nil, ErrReadOnly
	}

	f, err := ro.Open(name)
	if err != nil {
		return nil, err
	}

	return rofile{f}, nil
}

func (ro rofs) Remove(name string) error {
	return ErrReadOnly
}

func (ro rofs) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(ro.fs, fromAbs(name))
}

func (ro rofs) Mkdir(name string, perm fs.FileMode) error {
	return ErrReadOnly
}

func (ro rofs) MkdirAll(name string, perm fs.FileMode) error {
	return ErrReadOnly
}

func (ro rofs) RemoveAll(name string) error {
	return ErrReadOnly
}

type rofile struct{ fs.File }

func (f rofile) Write([]byte) (int, error) {
	return 0, ErrReadOnly
}
