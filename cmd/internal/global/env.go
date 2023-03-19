package global

import (
	_ "embed"

	"libdb.so/console/fs"
	"libdb.so/public"

	_ "libdb.so/console/programs/coreutils"
	_ "libdb.so/console/programs/hewwo"
	_ "libdb.so/console/programs/neofetch"
	_ "libdb.so/console/programs/spew"
	_ "libdb.so/console/programs/termio"
)

const InitialCwd = "/"

//go:embed shellrc
var RC string

var Filesystem = fs.ReadOnlyFS(public.FS)
