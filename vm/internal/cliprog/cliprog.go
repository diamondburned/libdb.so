package cliprog

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
	"libdb.so/vm"
)

// Wrap wraps a cli.App into a vm.Program.
func Wrap(app cli.App) vm.Program {
	app.UseShortOptionHandling = true
	app.Setup()
	if app.Name == "" {
		log.Println("cli: not registering app with empty name")
		return nil
	}
	return program{app}
}

type program struct {
	cli.App
}

func (p program) Name() string {
	return p.App.Name
}

func (p program) Usage() string {
	return p.App.Usage
}

func (p program) Run(ctx context.Context, env vm.Environment, args []string) error {
	p.Reader = env.Terminal.Stdin
	p.Writer = env.Terminal.Stdout
	p.ErrWriter = env.Terminal.Stderr
	p.ExitErrHandler = func(*cli.Context, error) {}
	return p.App.RunContext(ctx, args)
}
