//go:build js

package nsfw

import (
	"fmt"
	"io/fs"
	"strconv"
	"syscall/js"

	"libdb.so/vm/rwfs"
)

var (
	localStorage  = js.Global().Get("localStorage")
	dispatchEvent = js.Global().Get("dispatchEvent")
)

func IsEnabled() bool {
	nsfw := localStorage.Get("nsfw-v1")
	return nsfw.String() == "true"
}

func Enable() {
	setNSFW(true)
}

func Disable() {
	setNSFW(false)
}

func setNSFW(nsfw bool) {
	oldValue := localStorage.Get("nsfw-v1")
	newValue := js.ValueOf(strconv.FormatBool(nsfw))

	localStorage.Set("nsfw-v1", newValue)

	event := js.Global().Get("StorageEvent").New("storage", map[string]any{
		"key":      "nsfw-v1",
		"oldValue": oldValue,
		"newValue": newValue,
	})
	dispatchEvent.Invoke(event)
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
