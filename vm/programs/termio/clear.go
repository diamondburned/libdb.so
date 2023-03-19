package termio

import (
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(clear))
}

var clear = cli.App{
	Name:  "clear",
	Usage: "clear the terminal screen",
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)
		env.Print("\033[2J\033[1;1H")
		return nil
	},
}
