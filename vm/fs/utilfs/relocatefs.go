package utilfs

import (
	"io/fs"
	"path"
	"strings"
	"time"
)

type relocatedFS struct {
	fs      fs.FS
	newPath string
}

var _ fs.FS = relocatedFS{}

// RelocateFS returns a filesystem that relocates all paths in fs from oldPrefix
// to newPrefix.
func RelocateFS(newPath string) func(fs.FS) fs.FS {
	return func(fs fs.FS) fs.FS {
		return relocatedFS{fs, path.Clean(newPath)}
	}
}

func (r relocatedFS) Open(name string) (fs.File, error) {
	name = path.Clean(name)

	// Cases to be handled:
	//
	// - If newPath is "docs/files", and ReadDir is called with
	//   "docs", we have to generate a single DirEntry pointing to "files".
	// - If newPath is "docs/files", and ReadDir is called with "docs/files",
	//   then we pass through the entire listing.

	if strings.HasPrefix(name, r.newPath) {
		name = strings.TrimPrefix(name, r.newPath)
		name = strings.TrimPrefix(name, "/")
		name = path.Clean(name)

		return r.fs.Open(name)
	}

	if strings.HasPrefix(r.newPath, name) || name == "." {
		next := name
		name = path.Base(name)

		next = strings.TrimPrefix(r.newPath, next)
		next = strings.TrimPrefix(next, "/")
		next = path.Clean(next)
		next = popFirstPart(next)

		return relocatingDir{r: &r, name: name, next: next}, nil
	}

	return nil, fs.ErrNotExist
}

type relocatingDir struct {
	r    *relocatedFS
	name string
	next string
}

var (
	_ fs.File        = (*relocatingDir)(nil)
	_ fs.FileInfo    = (*relocatingDir)(nil)
	_ fs.DirEntry    = (*relocatingDir)(nil)
	_ fs.ReadDirFile = (*relocatingDir)(nil)
)

func (r relocatingDir) Name() string               { return r.name }
func (r relocatingDir) Size() int64                { return 0 }
func (r relocatingDir) Mode() fs.FileMode          { return fs.ModeDir | fs.ModePerm }
func (r relocatingDir) Type() fs.FileMode          { return fs.ModeDir }
func (r relocatingDir) IsDir() bool                { return true }
func (r relocatingDir) ModTime() time.Time         { return time.Time{} }
func (r relocatingDir) Sys() interface{}           { return nil }
func (r relocatingDir) Info() (fs.FileInfo, error) { return r, nil }
func (r relocatingDir) Stat() (fs.FileInfo, error) { return r, nil }
func (r relocatingDir) Read([]byte) (int, error)   { return 0, fs.ErrInvalid }
func (r relocatingDir) Close() error               { return nil }

func (r relocatingDir) ReadDir(n int) ([]fs.DirEntry, error) {
	return []fs.DirEntry{
		relocatingDir{
			r:    r.r,
			name: r.next,
		},
	}, nil
}

func popFirstPart(name string) string {
	first, _, _ := strings.Cut(name, "/")
	return first
}
