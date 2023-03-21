package kvfs

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"sync/atomic"
	"time"
)

// dmask is the default mask.
const dmask = 0777

type fsFileInfo struct {
	path string
	size int64
	time time.Time
	mode fs.FileMode
}

var _ fs.FileInfo = fsFileInfo{}

func dirInfo(store Store, path string, d StoredDirectory) fsFileInfo {
	return fsFileInfo{
		path: path,
		time: time.Unix(dirMtime(store, path, d.CreateTime), 0),
		mode: dmask | os.ModeDir,
	}
}

func fileInfo(path string, f StoredFile) fsFileInfo {
	return fsFileInfo{
		path: path,
		size: int64(len(f.Data)),
		time: time.Unix(f.ModTime, 0),
		mode: dmask,
	}
}

func (stat fsFileInfo) Name() string       { return path.Base(stat.path) }
func (stat fsFileInfo) Size() int64        { return stat.size }
func (stat fsFileInfo) ModTime() time.Time { return time.Time{} }
func (stat fsFileInfo) IsDir() bool        { return stat.mode&os.ModeDir != 0 }
func (stat fsFileInfo) Sys() any           { return stat }
func (stat fsFileInfo) Mode() fs.FileMode  { return stat.mode }

func (stat fsFileInfo) Type() fs.FileMode          { return stat.mode.Type() }
func (stat fsFileInfo) Info() (fs.FileInfo, error) { return stat, nil }

type fsFile struct {
	parent *FS
	info   fsFileInfo
	buf    bytes.Buffer
	flag   int // open mode
	closed int32
}

var _ fs.File = (*fsFile)(nil)

func (f *fsFile) Stat() (fs.FileInfo, error) {
	return f.info, nil
}

func (f *fsFile) Read(b []byte) (int, error) {
	if !flagRead(f.flag) {
		return 0, fs.ErrPermission
	}
	return f.buf.Read(b)
}

func (f *fsFile) Write(b []byte) (int, error) {
	if !flagWrite(f.flag) {
		return 0, fs.ErrPermission
	}
	return f.buf.Write(b)
}

func (f *fsFile) Close() error {
	if !atomic.CompareAndSwapInt32(&f.closed, 0, 1) {
		return fs.ErrClosed
	}
	if flagWrite(f.flag) {
		return f.parent.write(f)
	}
	return nil
}

type fsDir struct {
	parent *FS
	info   fsFileInfo
}

var _ fs.ReadDirFile = (*fsDir)(nil)

func (d *fsDir) Stat() (fs.FileInfo, error) { return d.info, nil }

func (d *fsDir) Read(b []byte) (int, error) { return 0, fs.ErrInvalid }

func (d *fsDir) Close() error { return nil }

func (d *fsDir) ReadDir(n int) ([]fs.DirEntry, error) {
	return d.parent.readDir(d.info.path, n)
}

func flagRead(flag int) bool {
	return flagHas(flag, os.O_RDONLY, os.O_RDWR)
}

func flagWrite(flag int) bool {
	return flagHas(flag, os.O_WRONLY, os.O_RDWR)
}

func flagHas(flag int, anyOf ...int) bool {
	for _, arg := range anyOf {
		if arg < 0b11 {
			// Since O_RDONLY starts at 0 (0b00) and O_WRONLY is 1 (0b01), we
			// have to AND the flag and compare equality.
			if flag&0b11 == arg {
				return true
			}
			continue
		}
		if flag&arg != 0 {
			return true
		}
	}
	return false
}
