package spew

import (
	"io/fs"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(app))
}

var app = cli.App{
	Name:  "spew",
	Usage: "spew file(s)",
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)
		log := vm.LoggerFromContext(c.Context)

		for _, arg := range c.Args().Slice() {
			path := path.Join(env.Cwd, arg)

			f, err := fs.ReadFile(env.Filesystem, path)
			if err != nil {
				log.Println("cat:", err)
				continue
			}

			spew.Fdump(env.Terminal.Stdout, f)
		}

		return nil
	},
}
