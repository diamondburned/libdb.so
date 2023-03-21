package rwfs

import (
	"io"
	"io/fs"
)

// FS implements a read-writable filesystem.
type FS interface {
	fs.FS
	// OpenFile opens the named file with specified flag. The returned file
	// may be read or written depending on the flag.
	OpenFile(name string, flag int, perm fs.FileMode) (File, error)
	// Remove removes the named file or (empty) directory.
	Remove(name string) error
	// Mkdir creates a new directory with the specified name and permission
	// bits.
	Mkdir(name string, perm fs.FileMode) error
	// MkdirAll creates a directory named path, along with any necessary
	// parents, and returns nil, or else returns an error.
	MkdirAll(name string, perm fs.FileMode) error
	// RemoveAll removes path and any children it contains. It removes
	// everything it can and returns errors using errors.Join.
	RemoveAll(name string) error
}

// File is a read-writable file.
type File interface {
	fs.File
	io.ReadWriteCloser
}

// Type aliases.
type (
	FileInfo = fs.FileInfo
	DirEntry = fs.DirEntry
	FileMode = fs.FileMode
)
