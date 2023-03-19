package global

import (
	_ "embed"

	_ "libdb.so/vm/programs/coreutils"
	_ "libdb.so/vm/programs/hewwo"
	_ "libdb.so/vm/programs/neofetch"
	_ "libdb.so/vm/programs/spew"
	_ "libdb.so/vm/programs/termio"
)

const InitialCwd = "/"

//go:embed shellrc
var RC string
