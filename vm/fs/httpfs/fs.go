package httpfs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
	"sync"

	"github.com/pkg/errors"
	"libdb.so/vm/fs/rwfs"
)

// FS is a file system that reads from an HTTP server.
type FS struct {
	root   FileTreeRoot
	client httpClient
}

var _ fs.FS = (*FS)(nil)

// New returns a new FS that obeys the given file tree. A cache may optionally
// be provided to cache file contents.
func New(client *http.Client, root FileTreeRoot, basePath string) fs.FS {
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

type unfetchedFS struct {
	getRoot func() *FS
	client  *http.Client
}

// NewFromURL returns a new FS that fetches the file tree from the given URL.
// The FS is silently fetched in the background.
func NewFromURL(client *http.Client, fsURL string) fs.FS {
	getRoot2 := func() (*FS, error) {
		resp, err := client.Get(fsURL)
		if err != nil {
			return nil, fmt.Errorf("fetching file tree: %w", err)
		}
		defer resp.Body.Close()

		var root FileTreeRoot
		if err := json.NewDecoder(resp.Body).Decode(&root); err != nil {
			return nil, fmt.Errorf("unmarshaling JSON file tree: %w", err)
		}

		return &FS{
			root: root,
			client: httpClient{
				client:   client,
				basePath: root.BaseURL,
			},
		}, nil
	}

	getRoot := func() *FS {
		fs, err := getRoot2()
		if err != nil {
			log.Println("fs warning:", err)
		}
		return fs
	}

	getRoot = sync.OnceValue(getRoot)
	go getRoot()

	return &unfetchedFS{
		client:  client,
		getRoot: getRoot,
	}
}

func (u *unfetchedFS) Open(path string) (fs.File, error) {
	f := u.getRoot()
	if f == nil {
		return nil, fs.ErrNotExist
	}
	return f.Open(path)
}
