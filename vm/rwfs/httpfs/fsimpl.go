package httpfs

import (
	"io"
	"io/fs"
	"sort"
	"time"
)

type fsFileInfo struct {
	name string
	size int64
	dir  bool
}

var _ fs.FileInfo = fsFileInfo{}

func dirInfo(name string) fsFileInfo {
	return fsFileInfo{
		name: name,
		size: 0,
		dir:  true,
	}
}

func fileInfo(name string, info FileInfo) fsFileInfo {
	return fsFileInfo{
		name: name,
		size: int64(info.Size),
	}
}

func (stat fsFileInfo) Name() string       { return stat.name }
func (stat fsFileInfo) Size() int64        { return stat.size }
func (stat fsFileInfo) ModTime() time.Time { return time.Time{} }
func (stat fsFileInfo) IsDir() bool        { return stat.dir }
func (stat fsFileInfo) Sys() any           { return stat }
func (stat fsFileInfo) Mode() fs.FileMode {
	mode := fs.FileMode(0444)
	if stat.IsDir() {
		mode |= 0111
		mode |= fs.ModeDir
	}
	return mode
}

func (stat fsFileInfo) Type() fs.FileMode          { return stat.Mode().Type() }
func (stat fsFileInfo) Info() (fs.FileInfo, error) { return stat, nil }

type fsFile struct {
	i fsFileInfo
	r io.ReadCloser
}

var _ fs.File = fsFile{}

func (f fsFile) Stat() (fs.FileInfo, error) {
	return f.i, nil
}

func (f fsFile) Read(b []byte) (int, error) {
	return f.r.Read(b)
}

func (f fsFile) Close() error {
	return f.r.Close()
}

type fsDir struct {
	i fsFileInfo
	d FileTree
}

var _ fs.ReadDirFile = fsDir{}

func (f fsDir) Stat() (fs.FileInfo, error) {
	return f.i, nil
}

func (f fsDir) Read(b []byte) (int, error) {
	return 0, fs.ErrInvalid
}

func (f fsDir) Close() error { return nil }

func (f fsDir) ReadDir(n int) ([]fs.DirEntry, error) {
	i := 0
	ents := make([]fs.DirEntry, 0, len(f.d))
	for name, df := range f.d {
		_, isDir := df.(FileTree)
		ents = append(ents, fsFileInfo{
			name: name,
			size: 0,
			dir:  isDir,
		})
		i++
		if n > 0 && i >= n {
			break
		}
	}
	sort.Slice(ents, func(i, j int) bool {
		return ents[i].Name() < ents[j].Name()
	})
	return ents, nil
}
