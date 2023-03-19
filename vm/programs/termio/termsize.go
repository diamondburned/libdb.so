package termio

import (
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(termsize))
}

var termsize = cli.App{
	Name:  "termsize",
	Usage: "print the terminal size",
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)
		q := env.Terminal.Query()
		env.Printf("%d %d", q.Width, q.Height)
		return nil
	},
}
