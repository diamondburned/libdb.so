package coreutils

import (
	"io"

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
			f, err := env.Open(arg)
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
