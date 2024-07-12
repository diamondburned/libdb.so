package rwfs

import (
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Split splits an absolute path into its components. The first element is
// always the directory within root, meaning it is not empty. If absPath is "/",
// then Split returns []string{}.
func Split(absPath string) []string {
	absPath = strings.TrimPrefix(path.Clean("/"+absPath), "/")
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

// Copy copies a file at srcPath to dstPath. dstPath and srcPath can be in
// different filesystems.
func Copy(dstFS FS, dstPath string, srcFS fs.FS, srcPath string) error {
	srcFile, err := srcFS.Open(srcPath)
	if err != nil {
		return &fs.PathError{
			Op:   "open",
			Path: srcPath,
			Err:  errors.Wrap(err, "failed to open source file"),
		}
	}
	defer srcFile.Close()

	srcStat, err := srcFile.Stat()
	if err != nil {
		return &fs.PathError{
			Op:   "stat",
			Path: srcPath,
			Err:  errors.Wrap(err, "failed to stat source file"),
		}
	}

	if srcStat.Size() == 0 {
		// There's nothing to copy. We're done.
		return nil
	}

	dstFile, err := dstFS.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, srcStat.Mode())
	if err != nil {
		return &fs.PathError{
			Op:   "open",
			Path: dstPath,
			Err:  errors.Wrap(err, "failed to create destination file"),
		}
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return &fs.PathError{
			Op:   "copy",
			Path: dstPath,
			Err:  errors.Wrap(err, "failed to copy file"),
		}
	}

	return nil
}

// DeduplicateDirEntries removes duplicate directory entries in-place. The first
// entry encountered is kept, and all subsequent entries are removed.
func DeduplicateDirEntries(dirEntries []fs.DirEntry) []fs.DirEntry {
	set := make(map[string]struct{}, len(dirEntries))
	deduped := dirEntries[:0]
	for _, dirEntry := range dirEntries {
		if _, ok := set[dirEntry.Name()]; ok {
			continue
		}
		set[dirEntry.Name()] = struct{}{}
		deduped = append(deduped, dirEntry)
	}
	return deduped
}
