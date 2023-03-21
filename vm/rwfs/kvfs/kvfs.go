package kvfs

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"libdb.so/vm/rwfs"
	"libdb.so/vm/rwfs/internal/fsutil"
)

// Store is a key-value store.
type Store interface {
	// Get returns the value for the given key.
	Get(string) (StoredValue, error)
	// Set sets the value for the given key.
	Set(string, StoredValue) error
	// Delete deletes the value for the given key.
	Delete(string) error
	// List returns a list of keys that match the given prefix. The given path
	// is guaranteed to end with a slash.
	List(prefix string, recursive bool) ([]PathedStoreValue, error)
}

// PathedStoreValue is specifically used to handle Store's List method.
type PathedStoreValue struct {
	StoredValue
	Path string
}

// StoredValue is a value stored in the store.
type StoredValue interface {
	storedValue()
}

func (StoredFile) storedValue()      {}
func (StoredDirectory) storedValue() {}

// UnmarshalStoredValue unmarshals the given stored value into the given value.
func UnmarshalStoredValue(b json.RawMessage) (StoredValue, error) {
	var isDir struct {
		IsDir bool `json:"isDir"`
	}

	if err := json.Unmarshal(b, &isDir); err != nil {
		return nil, errors.Wrap(err, "invalid json")
	}

	var v StoredValue
	var err error
	if isDir.IsDir {
		var d StoredDirectory
		err = json.Unmarshal(b, &d)
		v = d
	} else {
		var f StoredFile
		err = json.Unmarshal(b, &f)
		v = f
	}

	return v, err
}

// StoredFile is a file that is stored in a key-value store.
type StoredFile struct {
	// ModTime is the file's modification time.
	ModTime int64 `json:"mod_time"`
	// Data is the file's data.
	Data []byte `json:"data,omitempty"`
}

// StoredDirectory is a directory that is stored in a key-value store.
type StoredDirectory struct {
	// CreateTime is the directory's creation time.
	CreateTime int64 `json:"create_time"`
	// IsDir is always true.
	IsDir bool `json:"is_dir"`
}

// dirMtime returns the latest modification time of the given directory.
func dirMtime(store Store, dirpath string, createTime int64) int64 {
	time := createTime

	files, err := store.List(dirpath, false)
	if err != nil {
		return time
	}

	for _, file := range files {
		switch v := file.StoredValue.(type) {
		case StoredFile:
			if v.ModTime > time {
				time = v.ModTime
			}
		case StoredDirectory:
			if v.CreateTime > time {
				time = v.CreateTime
			}
		}
	}

	return time
}

// FS implements a key-value store-based file system.
type FS struct {
	store Store
	lock  sync.RWMutex
}

var (
	_ fs.FS        = (*FS)(nil)
	_ fs.ReadDirFS = (*FS)(nil)
	_ rwfs.FS      = (*FS)(nil)
)

// New returns a new FS that uses the given store. Be careful when constructing
// two different FS's with the same store, as they will conflict with each
// other.
func New(store Store) *FS {
	return &FS{store: store}
}

// Open implements fs.FS.
func (kvfs *FS) Open(fullpath string) (fs.File, error) {
	fullpath = clean(fullpath)

	kvfs.lock.RLock()
	defer kvfs.lock.RUnlock()

	v, err := kvfs.store.Get(fullpath)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "open",
			Path: fullpath,
			Err:  err,
		}
	}

	switch v := v.(type) {
	case StoredFile:
		return &fsFile{
			parent: kvfs,
			info:   fileInfo(fullpath, v),
			buf:    *bytes.NewBuffer(v.Data),
			flag:   os.O_RDONLY,
		}, nil
	case StoredDirectory:
		return &fsDir{
			parent: kvfs,
			info:   dirInfo(kvfs.store, fullpath, v),
		}, nil
	default:
		panic("unknown (impossible) stored value type")
	}
}

func (kvfs *FS) OpenFile(fullpath string, flag int, perm fs.FileMode) (rwfs.File, error) {
	fullpath = clean(fullpath)

	// Don't support O_RDWR because we don't have any way to do that.
	if flagHas(flag, os.O_RDWR) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: fullpath,
			Err:  errors.New("O_RDWR not supported"),
		}
	}

	if flagHas(flag, os.O_RDONLY) {
		kvfs.lock.RLock()
		defer kvfs.lock.RUnlock()
	} else {
		kvfs.lock.Lock()
		defer kvfs.lock.Unlock()
	}

	now := time.Now().Unix()
	var stored StoredFile

	v, err := kvfs.store.Get(fullpath)
	if err == nil {
		f, ok := v.(StoredFile)
		if ok {
			stored = f
		}
	}

	if flagHas(flag, os.O_CREATE) {
		// Update ModTime right now so we can reuse it.
		stored.ModTime = now

		if flagHas(flag, os.O_EXCL) && err == nil {
			// File exists but we want to only create the file when there's
			// none. Exit.
			return nil, &fs.PathError{
				Op:   "open",
				Path: fullpath,
				Err:  fs.ErrExist,
			}
		}
	} else {
		// No O_CREATE, so we must error if the file doesn't exist.
		if err != nil {
			return nil, &fs.PathError{
				Op:   "open",
				Path: fullpath,
				Err:  err,
			}
		}
	}

	if flagHas(flag, os.O_TRUNC, os.O_CREATE) {
		// Guarantee that the parent directory exists before we create this
		// file.
		if path.Dir(fullpath) != "/" {
			_, err := kvfs.store.Get(path.Dir(fullpath))
			if err != nil {
				return nil, &fs.PathError{
					Op:   "open",
					Path: fullpath,
					Err:  errors.Wrap(err, "failed to get parent directory"),
				}
			}
		}

		if !flagWrite(flag) { // user sanity check
			return nil, &fs.PathError{
				Op:   "open",
				Path: fullpath,
				Err:  errors.New("O_TRUNC requires write flag"),
			}
		}

		if flagHas(flag, os.O_TRUNC) {
			stored.Data = nil
		}

		if err := kvfs.store.Set(fullpath, stored); err != nil {
			return nil, &fs.PathError{
				Op:   "open",
				Path: fullpath,
				Err:  errors.Wrap(err, "failed to truncate file"),
			}
		}
	}

	var buf bytes.Buffer
	// We need the data in the buffer for reading. For writing, we're either
	// writing the data to the start or end (or overriding), so we don't need
	// to initialize the buffer.
	if flagRead(flag) {
		buf = *bytes.NewBuffer(stored.Data)
	}

	return &fsFile{
		parent: kvfs,
		info:   fileInfo(fullpath, stored),
		buf:    buf,
		flag:   flag,
	}, nil
}

func (kvfs *FS) write(file *fsFile) error {
	kvfs.lock.Lock()
	defer kvfs.lock.Unlock()

	now := time.Now().Unix()

	var f StoredFile
	var ok bool

	v, err := kvfs.store.Get(file.info.path)
	if err == nil {
		// Our existing file must've still been a file. Things might've changed
		// while the file was open, so we need to check again. We also have to
		// do this to get existing data for O_APPEND just in case.
		if f, ok = v.(StoredFile); !ok {
			return &fs.PathError{
				Op:   "write",
				Path: file.info.path,
				Err:  fs.ErrInvalid,
			}
		}
	}

	f.ModTime = now

	// Be careful with this. If we're appending, then we'll read again and
	// append to that file. This means that while the file is open, the appended
	// bytes might be very different.
	switch {
	case flagHas(file.flag, os.O_APPEND):
		f.Data = append(f.Data, file.buf.Bytes()...)
	case flagHas(file.flag, os.O_TRUNC):
		f.Data = file.buf.Bytes()
	default:
		f.Data = append(file.buf.Bytes(), f.Data...)
	}

	if err := kvfs.store.Set(file.info.path, f); err != nil {
		return &fs.PathError{
			Op:   "write",
			Path: file.info.path,
			Err:  err,
		}
	}

	return nil
}

func (rwfs *FS) ReadDir(fullpath string) ([]fs.DirEntry, error) {
	fullpath = clean(fullpath)
	return rwfs.readDir(fullpath, 0)
}

func (rwfs *FS) readDir(fullpath string, n int) ([]fs.DirEntry, error) {
	rwfs.lock.RLock()
	defer rwfs.lock.RUnlock()

	// Check that this is an actual directory.
	dir, err := rwfs.store.Get(fullpath)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "readdir",
			Path: fullpath,
			Err:  err,
		}
	}

	if _, ok := dir.(StoredDirectory); !ok {
		return nil, &fs.PathError{
			Op:   "readdir",
			Path: fullpath,
			Err:  fs.ErrInvalid,
		}
	}

	prefix := fullpath + "/"

	files, err := rwfs.store.List(prefix, false)
	if err != nil {
		return nil, err
	}

	if n > 0 && len(files) > n {
		files = files[:n]
	}

	entries := make([]fs.DirEntry, 0, len(files))
	for _, f := range files {
		name := strings.TrimPrefix(f.Path, prefix)
		// If the file has a slash, then it's in a subdirectory, so we'll ignore
		// it. Ideally, Store should do this, but we'll be lenient.
		if strings.Contains(name, "/") {
			continue
		}

		switch v := f.StoredValue.(type) {
		case StoredFile:
			entries = append(entries, fileInfo(f.Path, v))
		case StoredDirectory:
			entries = append(entries, dirInfo(rwfs.store, f.Path, v))
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil

}

func (rwfs *FS) Remove(fullpath string) error {
	rwfs.lock.Lock()
	defer rwfs.lock.Unlock()

	// Remove only supports removing files, not directories. We check for that.
	v, err := rwfs.store.Get(fullpath)
	if err != nil {
		// nothing to remove.
		return nil
	}

	_, ok := v.(StoredFile)
	if !ok {
		return &fs.PathError{
			Op:   "remove",
			Path: fullpath,
			Err:  fs.ErrInvalid,
		}
	}

	if err := rwfs.store.Delete(fullpath); err != nil {
		return &fs.PathError{
			Op:   "remove",
			Path: fullpath,
			Err:  err,
		}
	}

	return nil
}

func (rwfs *FS) Mkdir(fullpath string, perm fs.FileMode) error {
	fullpath = clean(fullpath)
	if fullpath == "/" {
		return nil
	}

	rwfs.lock.Lock()
	defer rwfs.lock.Unlock()

	now := time.Now().Unix()

	_, err := rwfs.store.Get(fullpath)
	if err == nil {
		return &fs.PathError{
			Op:   "mkdir",
			Path: fullpath,
			Err:  fs.ErrExist,
		}
	}

	value := StoredDirectory{
		CreateTime: now,
		IsDir:      true,
	}

	if err := rwfs.store.Set(fullpath, value); err != nil {
		return &fs.PathError{
			Op:   "mkdir",
			Path: fullpath,
			Err:  err,
		}
	}

	return nil
}

func (rwfs *FS) MkdirAll(fullpath string, perm fs.FileMode) error {
	fullpath = clean(fullpath)
	if fullpath == "/" {
		return nil
	}

	rwfs.lock.RLock()
	// Check if we already created this path.
	_, err := rwfs.store.Get(fullpath)
	rwfs.lock.RUnlock()

	if err == nil {
		return nil
	}

	// Slow path. If the path exists before we acquire this write lock, then our
	// loop still won't do anything, which is good.
	rwfs.lock.Lock()
	defer rwfs.lock.Unlock()

	now := time.Now().Unix()

	// We need to create the path. We'll do this by splitting the path and
	// creating each directory one by one.
	parts := fsutil.Split(fullpath)

	for i := range parts {
		path := "/" + strings.Join(parts[:i+1], "/")

		_, err := rwfs.store.Get(path)
		if err == nil {
			continue
		}

		if err := rwfs.store.Set(path, StoredDirectory{
			CreateTime: now,
			IsDir:      true,
		}); err != nil {
			return &fs.PathError{
				Op:   "mkdirall",
				Path: fullpath,
				Err:  err,
			}
		}
	}

	return nil
}

func (rwfs *FS) RemoveAll(fullpath string) error {
	fullpath = clean(fullpath)

	rwfs.lock.Lock()
	defer rwfs.lock.Unlock()

	files, err := rwfs.store.List(fullpath, true)
	if err == nil {
		// Path is a directory, remove everything it has.
		for _, file := range files {
			if err := rwfs.store.Delete(file.Path); err != nil {
				return &fs.PathError{
					Op:   "removeall",
					Path: fullpath,
					Err:  err,
				}
			}
		}
	}

	if fullpath == "/" {
		if err := rwfs.store.Delete(fullpath); err != nil {
			return &fs.PathError{
				Op:   "removeall",
				Path: fullpath,
				Err:  errors.Wrap(err, "failed to remove top-most directory"),
			}
		}
	}

	return nil
}

func clean(fullpath string) string {
	return path.Clean("/" + fullpath)
}
