package staticfs

import (
	"fmt"
	"testing"
	"time"

	stdfs "io/fs"

	"github.com/alecthomas/assert"
	"github.com/davecgh/go-spew/spew"
)

func TestUnflatten(t *testing.T) {
	tests := []struct {
		in  map[string]string
		out Directory
	}{
		{
			in: map[string]string{
				"/a/b/c": "hello",
				"/a/b/d": "world",
			},
			out: Directory{
				"a": Directory{
					"b": Directory{
						"c": File("hello"),
						"d": File("world"),
					},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			d, err := Unflatten(test.in)
			assert.NoError(t, err)
			assert.Equal(t, test.out, d)
		})
	}
}

func TestFS(t *testing.T) {
	fs := New(MustUnflatten(map[string]string{
		"/a/b/c": "hello",
		"/a/b/d": "world",
	}))

	t.Log(spew.Sdump(fs))

	root, err := fs.Open("/")
	assert.NoError(t, err)

	rootDir, ok := root.(stdfs.ReadDirFile)
	assert.True(t, ok)

	rootDirEntries, err := rootDir.ReadDir(-1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rootDirEntries))
	assert.Equal(t, "a", rootDirEntries[0].Name())

	b, err := fs.Open("/a/b")
	assert.NoError(t, err)

	bDir, ok := b.(stdfs.ReadDirFile)
	assert.True(t, ok)

	bDirEntries, err := bDir.ReadDir(-1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(bDirEntries))
	assert.Equal(t, "c", bDirEntries[0].Name())
	assert.Equal(t, "d", bDirEntries[1].Name())

	bstat, err := b.Stat()
	assert.NoError(t, err)
	assert.Equal(t, "b", bstat.Name())
	assert.Equal(t, int64(0), bstat.Size())
	assert.Equal(t, true, bstat.IsDir())
	assert.Equal(t, time.Time{}, bstat.ModTime())

	c, err := stdfs.ReadFile(fs, "/a/b/c")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(c))

	d, err := stdfs.ReadFile(fs, "/a/b/d")
	assert.NoError(t, err)
	assert.Equal(t, "world", string(d))
}
