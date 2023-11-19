//go:build js

package nsfw

import (
	"fmt"
	"io/fs"

	"libdb.so/vm/internal/vars"
	"libdb.so/vm/rwfs"
)

var Variable = vars.
	New[bool]("nsfw-v1").
	WithDefault(false).
	WithHidden(true)

func IsEnabled() bool {
	return Variable.Getz()
}

func Enable() {
	Variable.Set(true)
}

func Disable() {
	Variable.Set(false)
}

func WrapFS(rofs fs.FS) fs.FS {
	return nsfwFS{ro: rofs}
}

type nsfwFS struct {
	ro fs.FS
}

var errDenied = fmt.Errorf("%w (see `help')", fs.ErrPermission)

func (f nsfwFS) Open(path string) (fs.File, error) {
	if pathHasNSFW(path) && !IsEnabled() {
		return nil, errDenied
	}
	return f.ro.Open(path)
}

func (f nsfwFS) Stat(path string) (fs.FileInfo, error) {
	if pathHasNSFW(path) && !IsEnabled() {
		return nil, errDenied
	}
	return fs.Stat(f.ro, path)
}

func pathHasNSFW(path string) bool {
	for _, part := range rwfs.Split(path) {
		if part == ".nsfw" {
			return true
		}
	}

	return false
}
