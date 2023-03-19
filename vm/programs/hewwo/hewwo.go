package hewwo

import (
	"context"

	"libdb.so/vm"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(program{})
}

type program struct{}

func (p program) Name() string {
	return "hewwo"
}

func (p program) Run(ctx context.Context, env *vm.Environment, args []string) error {
	if len(args) != 1 {
		return &vm.UsageError{Usage: "hewwo"}
	}
	env.Println("hewwo go!")
	return nil
}
