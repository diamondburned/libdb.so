package coreutils

import (
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(sleep))
}

var sleep = cli.App{
	Name:      "sleep",
	Usage:     "sleep for NUMBER seconds or a given duration",
	UsageText: `sleep NUMBER[SUFFIX]`,
	Action: func(c *cli.Context) error {
		if c.NArg() != 1 {
			return &vm.UsageError{Usage: "sleep NUMBER[SUFFIX]"}
		}

		d, err := time.ParseDuration(c.Args().First())
		if err != nil {
			i, ierr := strconv.ParseFloat(c.Args().First(), 64)
			if ierr != nil {
				return errors.Wrap(err, "invalid duration or number")
			}
			d = time.Duration(math.Ceil(i * float64(time.Second)))
		}

		timer := time.NewTimer(d)
		defer timer.Stop()

		select {
		case <-timer.C:
			return nil
		case <-c.Context.Done():
			return nil
		}
	},
}
