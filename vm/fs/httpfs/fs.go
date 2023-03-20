package httpfs

import (
	"io/fs"
	"net/http"

	"github.com/pkg/errors"
	"libdb.so/vm/fs/internal/fsutil"
)

// FS is a file system that reads from an HTTP server.
type FS struct {
	tree   FileTree
	client httpClient
}

var _ fs.FS = (*FS)(nil)

// New returns a new FS that obeys the given file tree. A cache may optionally
// be provided to cache file contents.
func New(client http.Client, tree FileTree, basePath string) *FS {
	return &FS{
		tree: tree,
		client: httpClient{
			client:   client,
			basePath: basePath,
		},
	}
}

func (h *FS) Open(path string) (fs.File, error) {
	parts := fsutil.Split(path)
	if len(parts) == 0 {
		return fsDir{
			i: dirInfo(path),
			d: h.tree,
		}, nil
	}

	current := h.tree
	for i, part := range parts {
		entry, ok := current[part]
		if !ok {
			return nil, fs.ErrNotExist
		}

		if i == len(parts)-1 {
			switch entry := entry.(type) {
			case FileInfo:
				r, err := h.client.get(fsutil.JoinAbs(parts), entry)
				if err != nil {
					return nil, errors.Wrap(err, "httpfs: error getting file")
				}
				return fsFile{
					i: fileInfo(part, entry),
					r: r,
				}, nil
			case FileTree:
				return fsDir{
					i: dirInfo(part),
					d: entry,
				}, nil
			}
		}

		dir, ok := entry.(FileTree)
		if !ok {
			return nil, fs.ErrNotExist
		}

		current = dir
	}

	return nil, fs.ErrNotExist
}
