package fsutil

import (
	"path"
	"strings"
)

// Split splits an absolute path into its components. The first element is
// always the directory within root, meaning it is not empty. If absPath is "/",
// then Split returns []string{}.
func Split(absPath string) []string {
	absPath = strings.TrimPrefix(path.Clean(absPath), "/")
	if absPath == "" {
		return []string{}
	}

	return strings.Split(absPath, "/")
}

// JoinAbs joins a path to an absolute path.
func JoinAbs(parts []string) string {
	return "/" + path.Join(parts...)
}

// ConvertAbs converts an absolute path to an io/fs-compatible path. The
// returned path will still be absolute, but it will not start with a
// leading slash. If the path is the root, "." is returned.
func ConvertAbs(abs string) string {
	cleaned := path.Clean(abs)
	if !strings.HasPrefix(cleaned, "/") {
		return cleaned
	}

	cleaned = strings.TrimPrefix(cleaned, "/")
	if cleaned == "" {
		cleaned = "."
	}

	return cleaned
}
