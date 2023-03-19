// Package public contains contents in the public/ folder.
package public

import "embed"

//go:embed *
var FS embed.FS
