package httpfs

import (
	"encoding/json"
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
	Size int64  `json:"size"`
	Hash []byte `json:"hash"`
}

// FileTree is a tree of files and directories. Directories will have its keys
// end with a slash, while files will not.
type FileTree map[string]FileTreeValue

func (t *FileTree) UnmarshalJSON(b []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	*t = make(FileTree)

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
