package coreutils

import (
	"io"

	"github.com/urfave/cli/v3"
	"libdb.so/console"
	"libdb.so/console/internal/cliprog"
	"libdb.so/console/programs"
)

func init() {
	programs.Register(cliprog.Wrap(cat))
}

var cat = cli.App{
	Name:  "cat",
	Usage: "concatenate files and print on the standard output",
	Action: func(c *cli.Context) error {
		env := console.EnvironmentFromContext(c.Context)
		log := console.LoggerFromContext(c.Context)

		for _, arg := range c.Args().Slice() {
			f, err := env.Filesystem.Open(arg)
			if err != nil {
				log.Println("cat:", err)
				continue
			}

			if _, err = io.Copy(env.Terminal.Stdout, f); err != nil {
				log.Println("cat: io.Copy:", err)
			}

			f.Close()
		}

		return nil
	},
}
