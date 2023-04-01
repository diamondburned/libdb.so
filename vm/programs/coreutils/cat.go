package coreutils

import (
	"io"
	"path"

	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(cat))
}

var cat = cli.App{
	Name:      "cat",
	Usage:     "concatenate files and print on the standard output",
	UsageText: `cat [FILE]...`,
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)
		log := vm.LoggerFromContext(c.Context)

		for _, arg := range c.Args().Slice() {
			path := path.Join(env.Cwd, arg)

			f, err := env.Open(path)
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
