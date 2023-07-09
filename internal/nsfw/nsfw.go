//go:build js

package nsfw

import (
	"fmt"
	"io/fs"
	"syscall/js"

	"libdb.so/vm/rwfs"
)

func IsEnabled() bool {
	local := js.Global().Get("localStorage")
	nsfw := local.Get("nsfw-v1")
	return nsfw.String() == "true"
}

func Enable() {
	local := js.Global().Get("localStorage")
	local.Set("nsfw-v1", true)
}

func Disable() {
	local := js.Global().Get("localStorage")
	local.Set("nsfw-v1", false)
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
