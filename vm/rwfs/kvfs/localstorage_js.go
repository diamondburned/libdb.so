package kvfs

import (
	"encoding/json"
	"io/fs"
	"path"
	"strings"
	"syscall/js"

	"github.com/pkg/errors"
)

const localStoragePrefix = "__kvfs"

// LocalStorage is a simple key-value store that is persisted to the
// browser's local storage.
func LocalStorage() Store {
	return localStorage{
		js.Global().Get("localStorage"),
		js.Global().Get("Object"),
	}
}

var _ Store = localStorage{}

type localStorage struct {
	js     js.Value
	object js.Value
}

func (s localStorage) Get(fullpath string) (StoredValue, error) {
	fullpath = path.Clean("/" + fullpath)
	key := localStoragePrefix + fullpath

	jsv := s.js.Call("getItem", key)
	if jsv.IsNull() {
		return nil, fs.ErrNotExist
	}

	v, err := UnmarshalStoredValue([]byte(jsv.String()))
	if err != nil {
		return nil, errors.Wrap(err, "file contains invalid json")
	}

	return v, nil
}

func (s localStorage) Set(fullpath string, v StoredValue) error {
	fullpath = path.Clean("/" + fullpath)
	key := localStoragePrefix + fullpath

	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "failed to marshal value")
	}

	s.js.Call("setItem", key, string(b))

	return nil
}

func (s localStorage) Delete(fullpath string) error {
	fullpath = path.Clean("/" + fullpath)
	key := localStoragePrefix + fullpath

	s.js.Call("removeItem", key)
	return nil
}

func (s localStorage) List(fullpath string, recursive bool) ([]PathedStoreValue, error) {
	// We assume fullpath is already a directory.
	fullpath = path.Clean("/"+fullpath) + "/"
	key := localStoragePrefix + fullpath

	var values []PathedStoreValue
	for i := 0; i < s.js.Get("length").Int(); i++ {
		k := s.js.Call("key", i).String()
		if !strings.HasPrefix(k, key) {
			continue
		}

		filepath := strings.TrimPrefix(k, localStoragePrefix)
		filename := strings.TrimPrefix(filepath, fullpath)
		if !recursive && strings.Contains(filename, "/") {
			// Not a direct child. Skip.
			continue
		}

		jsv := s.js.Call("getItem", k)
		if jsv.IsNull() {
			// Shouldn't happen, but whatever.
			continue
		}

		v, err := UnmarshalStoredValue([]byte(jsv.String()))
		if err != nil {
			return nil, errors.Wrap(err, "file contains invalid json")
		}

		values = append(values, PathedStoreValue{
			StoredValue: v,
			Path:        filepath,
		})
	}

	return values, nil
}
