package staticfs

import (
	"sort"
	"strings"
	"time"

	stdfs "io/fs"
)

type fsFileInfo struct {
	name string
	size int64
	dir  bool
}

var _ stdfs.FileInfo = fsFileInfo{}

func dirInfo(name string, d Directory) fsFileInfo {
	return fsFileInfo{
		name: name,
		size: 0,
		dir:  true,
	}
}

func fileInfo(name string, f File) fsFileInfo {
	return fsFileInfo{
		name: name,
		size: int64(len(f)),
	}
}

func (stat fsFileInfo) Name() string       { return stat.name }
func (stat fsFileInfo) Size() int64        { return stat.size }
func (stat fsFileInfo) ModTime() time.Time { return time.Time{} }
func (stat fsFileInfo) IsDir() bool        { return stat.dir }
func (stat fsFileInfo) Sys() any           { return stat }
func (stat fsFileInfo) Mode() stdfs.FileMode {
	if stat.IsDir() {
		return 1444
	}
	return 0444
}

func (stat fsFileInfo) Type() stdfs.FileMode          { return stat.Mode().Type() }
func (stat fsFileInfo) Info() (stdfs.FileInfo, error) { return stat, nil }

type fsFile struct {
	i fsFileInfo
	r *strings.Reader
}

var _ stdfs.File = fsFile{}

func (f fsFile) Stat() (stdfs.FileInfo, error) {
	return f.i, nil
}

func (f fsFile) Read(b []byte) (int, error) {
	return f.r.Read(b)
}

func (f fsFile) Close() error { return nil }

type fsDir struct {
	i fsFileInfo
	d Directory
}

var _ stdfs.ReadDirFile = fsDir{}

func (f fsDir) Stat() (stdfs.FileInfo, error) {
	return f.i, nil
}

func (f fsDir) Read(b []byte) (int, error) {
	return 0, stdfs.ErrInvalid
}

func (f fsDir) Close() error { return nil }

func (f fsDir) ReadDir(n int) ([]stdfs.DirEntry, error) {
	i := 0
	ents := make([]stdfs.DirEntry, 0, len(f.d))
	for name, df := range f.d {
		_, isDir := df.(Directory)
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
