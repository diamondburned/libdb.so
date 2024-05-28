package httpfs

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"libdb.so/vm/rwfs"
)

// FS is a file system that reads from an HTTP server.
type FS struct {
	root   FileTreeRoot
	client httpClient
}

var _ fs.FS = (*FS)(nil)

// New returns a new FS that obeys the given file tree. A cache may optionally
// be provided to cache file contents.
func New(client http.Client, root FileTreeRoot, basePath string) *FS {
	return &FS{
		root: root,
		client: httpClient{
			client:   client,
			basePath: path.Join(root.BaseURL, basePath),
		},
	}
}

func (h *FS) Open(path string) (fs.File, error) {
	parts := rwfs.Split(path)
	if len(parts) == 0 {
		return fsDir{
			i: dirInfo(path),
			d: h.root.Tree,
		}, nil
	}

	current := h.root.Tree
	for i, part := range parts {
		entry, ok := current[part]
		if !ok {
			return nil, fs.ErrNotExist
		}

		if i == len(parts)-1 {
			switch entry := entry.(type) {
			case FileInfo:
				r, err := h.client.get(rwfs.JoinAbs(parts), entry)
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

			break
		}

		dir, ok := entry.(FileTree)
		if !ok {
			return nil, fs.ErrNotExist
		}

		current = dir
	}

	return nil, fs.ErrNotExist
}
