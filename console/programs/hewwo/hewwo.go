package hewwo

import (
	"fmt"

	"libdb.so/console"
	"libdb.so/console/programs"
)

func init() {
	programs.Register(program{})
}

type program struct{}

func (p program) Name() string { return "hewwo" }

func (p program) Run(term *console.Terminal, args []string) error {
	if len(args) != 1 {
		return &console.UsageError{Usage: "hewwo"}
	}

	fmt.Fprintln(term, "hewwo go!")
	return nil
}
