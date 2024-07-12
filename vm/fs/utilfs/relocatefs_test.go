package utilfs

import (
	"io"
	"io/fs"
	"testing"

	"github.com/alecthomas/assert/v2"
	"libdb.so/vm/fs/kvfs"
)

func TestRelocateFS(t *testing.T) {
	store := kvfs.MemoryStorageFromExisting(map[string]kvfs.StoredValue{
		"/d":   kvfs.StoredDirectory{IsDir: true},
		"/d/e": kvfs.StoredFile{Data: []byte("hello")},
	})
	testFS := kvfs.New(store)

	relocatedFS := Apply(testFS, RelocateFS("a/b/c"))

	t.Run("a/b/c/d/e", func(t *testing.T) {
		f, err := relocatedFS.Open("a/b/c/d/e")
		assert.NoError(t, err)

		b, err := io.ReadAll(f)
		f.Close()
		assert.NoError(t, err)
		assert.Equal(t, "hello", string(b), "file contents")
	})

	t.Run("a/b/c/d", func(t *testing.T) {
		files, err := fs.ReadDir(relocatedFS, "a/b/c/d")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(files), "number of files")
		assert.Equal(t, "e", files[0].Name(), "name")
	})

	t.Run("a/b/c", func(t *testing.T) {
		files, err := fs.ReadDir(relocatedFS, "a/b/c")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(files), "number of files")
		assert.Equal(t, "d", files[0].Name(), "name")
	})

	t.Run("a/b", func(t *testing.T) {
		files, err := fs.ReadDir(relocatedFS, "a/b")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(files), "number of files")
		assert.Equal(t, "c", files[0].Name(), "name")
	})

	t.Run("a", func(t *testing.T) {
		files, err := fs.ReadDir(relocatedFS, "a")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(files), "number of files")
		assert.Equal(t, "b", files[0].Name(), "name")
	})

	t.Run("-", func(t *testing.T) {
		files, err := fs.ReadDir(relocatedFS, "")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(files), "number of files")
		assert.Equal(t, "a", files[0].Name(), "name")
	})
}
