package cliprog

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
	"libdb.so/console"
)

// Wrap wraps a cli.App into a console.Program.
func Wrap(app cli.App) console.Program {
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

func (p program) Run(ctx context.Context, env *console.Environment, args []string) error {
	p.Reader = env.Terminal.Stdin
	p.Writer = env.Terminal.Stdout
	p.ErrWriter = env.Terminal.Stderr
	p.ExitErrHandler = func(*cli.Context, error) {}
	return p.App.RunContext(ctx, args)
}
