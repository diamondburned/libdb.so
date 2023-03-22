package kvfs

import (
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/davecgh/go-spew/spew"
)

func TestFS(t *testing.T) {
	store := MemoryStorage().(*memoryStorage)
	rwfs := New(store)

	t.Run("erroneous_create", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		var err error

		_, err = rwfs.OpenFile("foo/bar/baz", os.O_CREATE, 0)
		assert.Error(t, err)

		_, err = rwfs.OpenFile("foo/bar", os.O_CREATE, 0)
		assert.Error(t, err)
	})

	t.Run("mkdirAll", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		var err error

		err = rwfs.MkdirAll("foo/bar", 0)
		assert.NoError(t, err)

		err = rwfs.MkdirAll("foo/bar", 0)
		assert.NoError(t, err)

		root, err := rwfs.ReadDir("/")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(root))
		assert.Equal(t, "foo", root[0].Name())
	})

	t.Run("create", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		f, err := rwfs.OpenFile("foo/bar/baz", os.O_WRONLY|os.O_CREATE, 0)
		assert.NoError(t, err)

		defer f.Close()

		_, err = f.Write([]byte("hello"))
		assert.NoError(t, err)

		assert.NoError(t, f.Close())

		stat, err := f.Stat()
		assert.NoError(t, err)
		assert.Equal(t, "baz", stat.Name())
		assert.Equal(t, false, stat.IsDir())
		assert.Equal(t, 0777, stat.Mode())

		impl, ok := f.(*fsFile)
		assert.True(t, ok)
		assert.True(t, flagHas(impl.flag, os.O_WRONLY), "must have WRONLY flag")
		assert.Equal(t, rwfs, impl.parent)
		assert.Equal(t, 0, impl.info.mode&os.ModeDir)
		assert.Equal(t, 1, impl.closed, "must be closed")
		// Internally, we store paths with a prefixing slash just to make the
		// implementation easier.
		assert.Equal(t, "/foo/bar/baz", impl.info.path)
	})

	t.Run("readDirFile", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		f, err := rwfs.Open("foo/bar")
		assert.NoError(t, err)

		defer f.Close()

		readDirFile, ok := f.(fs.ReadDirFile)
		assert.True(t, ok)

		entries, err := readDirFile.ReadDir(0)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(entries))
		assert.False(t, entries[0].IsDir(), "must be a file")
		assert.Equal(t, "baz", entries[0].Name())
		assert.Equal(t, 0, entries[0].Type())
	})

	t.Run("readDirFS", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		entries, err := rwfs.ReadDir("foo/bar")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(entries))
		assert.False(t, entries[0].IsDir(), "must be a file")
		assert.Equal(t, "baz", entries[0].Name())
		assert.Equal(t, 0, entries[0].Type())
	})

	t.Run("read", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		test := func(f fs.File) {
			b, err := io.ReadAll(f)
			assert.NoError(t, err)
			assert.Equal(t, "hello", string(b))

			impl, ok := f.(*fsFile)
			assert.True(t, ok)
			assert.Equal(t, rwfs, impl.parent)
			// Internal is different from public as written above.
			assert.Equal(t, "/foo/bar/baz", impl.info.path)
		}

		f, err := rwfs.Open("foo/bar/baz")
		assert.NoError(t, err)
		defer f.Close()
		test(f)

		f, err = rwfs.OpenFile("foo/bar/baz", os.O_RDONLY, 0)
		assert.NoError(t, err)
		defer f.Close()
		test(f)
	})

	t.Run("append", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		f, err := rwfs.OpenFile("foo/bar/baz", os.O_WRONLY|os.O_APPEND, 0)
		assert.NoError(t, err)

		defer f.Close()

		_, err = f.Write([]byte(" yadda yadda"))
		assert.NoError(t, err)

		assert.NoError(t, f.Close())
	})

	t.Run("read-appended", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		f, err := rwfs.Open("foo/bar/baz")
		assert.NoError(t, err)
		defer f.Close()

		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, "hello yadda yadda", string(b))
	})

	t.Run("truncate", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		f, err := rwfs.OpenFile("foo/bar/baz", os.O_WRONLY|os.O_TRUNC, 0)
		assert.NoError(t, err)

		defer f.Close()

		_, err = f.Write([]byte("goodbye"))
		assert.NoError(t, err)

		assert.NoError(t, f.Close())
	})

	t.Run("read-truncated", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		f, err := rwfs.OpenFile("foo/bar/baz", os.O_RDONLY, 0)
		assert.NoError(t, err)

		defer f.Close()

		b, err := io.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, "goodbye", string(b))
	})

	t.Run("delete", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		err := rwfs.Remove("foo/bar/baz")
		assert.NoError(t, err)

		_, err = rwfs.OpenFile("foo/bar/baz", os.O_RDONLY, 0)
		assert.Error(t, err)

		list, err := fs.ReadDir(rwfs, "foo/bar")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(list), "directory should be empty:\n%s", spew.Sdump(list))
	})

	t.Run("delete dir", func(t *testing.T) {
		t.Cleanup(func() { t.Log(spew.Sdump(store.m)) })

		err := rwfs.RemoveAll("foo")
		assert.NoError(t, err)

		_, err = fs.ReadDir(rwfs, "foo/bar")
		assert.Error(t, err)

		_, err = fs.ReadDir(rwfs, "foo")
		assert.Error(t, err)
	})

	t.Log(spew.Sdump(store.m))
}
