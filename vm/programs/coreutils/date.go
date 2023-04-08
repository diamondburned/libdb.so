package coreutils

import (
	"errors"
	"fmt"
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
		&cli.StringFlag{
			Name:    "date",
			Aliases: []string{"d"},
			Usage:   "display time described by STRING, not 'now'; requires -f",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "use FORMAT for date parsing; see godate -h",
		},
	},
	Action: func(c *cli.Context) error {
		switch c.Args().First() {
		case "diamond", "me":
			fmt.Println("I love you <3")
			return nil
		}

		env := vm.EnvironmentFromContext(c.Context)

		now := time.Now()
		if c.Bool("date") {
			if !c.Bool("format") {
				return errors.New("date: -d requires -f")
			}

			t, err := time.Parse(c.String("format"), c.String("date"))
			if err != nil {
				return fmt.Errorf("date: invalid time %q: %w", c.String("date"), err)
			}

			now = t
		}

		if c.Bool("iso-8601") {
			env.Println(now.Format(time.RFC3339))
		} else {
			env.Println(now.Format(time.UnixDate))
		}

		return nil
	},
}

var godate = cli.App{
	Name:      "godate",
	Usage:     "convert strftime to Go time format",
	UsageText: `godate <strftime format>`,
	// TODO
}
