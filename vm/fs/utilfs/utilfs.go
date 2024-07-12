package utilfs

import "io/fs"

// Apply applies a series of appliers to a filesystem.
func Apply(fs fs.FS, appliers ...func(fs.FS) fs.FS) fs.FS {
	for _, apply := range appliers {
		fs = apply(fs)
	}
	return fs
}
