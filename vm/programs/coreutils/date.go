package coreutils

import (
	"time"

	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(date))
}

var date = cli.App{
	Name:  "date",
	Usage: "print the local date and time",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "iso-8601",
			Aliases: []string{"I", "rfc-3339"},
		},
	},
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)
		if c.Bool("iso-8601") {
			env.Println(time.Now().Format(time.RFC3339))
		} else {
			env.Println(time.Now().Format(time.UnixDate))
		}
		return nil
	},
}
