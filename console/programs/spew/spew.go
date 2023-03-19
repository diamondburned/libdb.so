package spew

import (
	"io/fs"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v3"
	"libdb.so/console"
	"libdb.so/console/internal/cliprog"
	"libdb.so/console/programs"
)

func init() {
	programs.Register(cliprog.Wrap(app))
}

var app = cli.App{
	Name:  "spew",
	Usage: "spew file(s)",
	Action: func(c *cli.Context) error {
		env := console.EnvironmentFromContext(c.Context)
		log := console.LoggerFromContext(c.Context)

		for _, arg := range c.Args().Slice() {
			f, err := fs.ReadFile(env.Filesystem, arg)
			if err != nil {
				log.Println("cat:", err)
				continue
			}

			spew.Fdump(env.Terminal.Stdout, f)
		}

		return nil
	},
}
