package hewwo

import (
	"context"

	"libdb.so/console"
	"libdb.so/console/programs"
)

func init() {
	programs.Register(program{})
}

type program struct{}

func (p program) Name() string {
	return "hewwo"
}

func (p program) Run(ctx context.Context, env *console.Environment, args []string) error {
	if len(args) != 1 {
		return &console.UsageError{Usage: "hewwo"}
	}
	env.Println("hewwo go!")
	return nil
}
