package staticfs

import (
	"fmt"
	"path"
	"strings"

	stdfs "io/fs"

	"libdb.so/vm/fs/internal/fsutil"
)

// Directory is a read-only directory constructed from a map of file contents.
type Directory map[string]DirEntry

// Unflatten returns a new filesystem constructed from a map of file paths to
// contents.
func Unflatten(files map[string]string) (Directory, error) {
	dir := make(Directory)

	for fpath, contents := range files {
		if !path.IsAbs(fpath) {
			return nil, fmt.Errorf("path %q is not absolute", fpath)
		}

		parts := strings.Split(path.Clean(fpath), "/")
		parts = parts[1:]
		current := dir

		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = File(contents)
				break
			}

			if _, ok := current[part]; !ok {
				current[part] = make(Directory)
			}

			current = current[part].(Directory)
		}
	}

	return dir, nil
}

// MustUnflatten is like Unflatten, but panics if the map of files is invalid.
func MustUnflatten(files map[string]string) Directory {
	dir, err := Unflatten(files)
	if err != nil {
		panic(err)
	}

	return dir
}

func (d Directory) at(parts []string) (DirEntry, string, error) {
	current := d

	for i, part := range parts {
		entry, ok := current[part]
		if !ok {
			return nil, "", stdfs.ErrNotExist
		}

		if i == len(parts)-1 {
			return entry, part, nil
		}

		dir, ok := entry.(Directory)
		if !ok {
			return nil, "", stdfs.ErrNotExist
		}

		current = dir
	}

	return nil, "", stdfs.ErrNotExist
}

// File is a file in a static filesystem.
type File string

// DirEntry is either a File or a Directory.
type DirEntry interface {
	dirEntry()
}

func (File) dirEntry()      {}
func (Directory) dirEntry() {}

// FS is a read-only filesystem constructed from a map of file contents.
type FS Directory

// New returns a new filesystem constructed from a map of file contents.
func New(dir Directory) FS {
	return FS(dir)
}

// Open returns a file at the given path.
func (fs FS) Open(path string) (stdfs.File, error) {
	entry, name, err := fs.at(path)
	if err != nil {
		return nil, err
	}

	switch entry := entry.(type) {
	case File:
		return fsFile{
			i: fileInfo(name, entry),
			r: strings.NewReader(string(entry)),
		}, nil
	case Directory:
		return fsDir{
			i: dirInfo(name, entry),
			d: entry,
		}, nil
	default:
		panic("unreachable")
	}
}

func (fs FS) at(fpath string) (DirEntry, string, error) {
	parts := fsutil.Split(fpath)
	if len(parts) == 0 {
		return Directory(fs), ".", nil
	}

	return Directory(fs).at(parts)

	// if fpath == "/" {
	// 	return Directory(fs), ".", nil
	// }
	// parts := strings.Split(path.Clean(fpath), "/")
	// if len(parts) > 0 && parts[0] == "" {
	// 	parts = parts[1:]
	// }
	// return Directory(fs).at(parts)
}
