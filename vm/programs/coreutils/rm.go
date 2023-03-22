package coreutils

import (
	"errors"

	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(rm))
}

var rm = cli.App{
	Name:      "rm",
	Usage:     "remove files or directories",
	UsageText: `rm [OPTION]... FILE...`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "remove directories and their contents recursively",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "no-op, for compatibility only",
		},
	},
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)
		log := vm.LoggerFromContext(c.Context)

		rm := env.Filesystem.Remove
		if c.Bool("recursive") {
			rm = env.Filesystem.RemoveAll
		}

		var failed bool
		for _, arg := range c.Args().Slice() {
			if err := rm(arg); err != nil {
				log.Println("rm:", err)
				failed = true
			}
		}

		if failed {
			return errors.New("failed to remove one or more files")
		}

		return nil
	},
}
