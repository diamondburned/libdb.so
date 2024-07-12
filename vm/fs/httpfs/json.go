package httpfs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FileTreeValue is either a FileTree or a FileSize.
type FileTreeValue interface {
	fileTreeValue()
}

func (FileInfo) fileTreeValue() {}
func (FileTree) fileTreeValue() {}

// FileInfo is the information of a file.
type FileInfo struct {
	Size int64 `json:"size"`
}

// FileTree is a tree of files and directories. Directories will have its keys
// end with a slash, while files will not.
type FileTree map[string]FileTreeValue

func (t *FileTree) UnmarshalJSON(b []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	*t = make(map[string]FileTreeValue, len(m))
	for k, b := range m {
		if strings.HasSuffix(k, "/") {
			var v FileTree
			if err := json.Unmarshal(b, &v); err != nil {
				return err
			}
			k = strings.TrimSuffix(k, "/")
			(*t)[k] = v
		} else {
			var v FileInfo
			if err := json.Unmarshal(b, &v); err != nil {
				return err
			}
			(*t)[k] = v
		}
	}

	return nil
}

// FileTreeRoot is a root file tree.
type FileTreeRoot struct {
	// BaseURL is the base URL of the file tree. Iti is prepended to the file
	// paths when fetching files.
	BaseURL string `json:"base_url"`
	// Tree is the file tree.
	Tree FileTree `json:"tree"`
}

func (t *FileTreeRoot) UnmarshalJSON(b []byte) error {
	var versionedTree struct {
		Version int `json:"$version"`
	}

	if err := json.Unmarshal(b, &versionedTree); err != nil {
		return fmt.Errorf("unmarshal file tree version: %w", err)
	}

	switch versionedTree.Version {
	case 0:
		return json.Unmarshal(b, &t.Tree)
	case 1:
		return fmt.Errorf("file tree version 1 is not used")
	case 2:
		type RawFileTreeRoot FileTreeRoot
		return json.Unmarshal(b, (*RawFileTreeRoot)(t))
	default:
		return fmt.Errorf("unknown file tree version: %d", versionedTree.Version)
	}
}
