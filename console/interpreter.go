package console

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	stderrors "errors"
	stdfs "io/fs"

	"github.com/pkg/errors"
	"libdb.so/console/fs"
	"libdb.so/console/internal/liner"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Interpreter is the main interpreter for the larger program. It manages
// prompting the user, printing to console, and running programs.
type Interpreter struct {
	shParser *syntax.Parser
	shRunner *interp.Runner
	prompter *liner.State
	logger   *log.Logger
	env      *Environment
	closes   []func() error
	opts     InterpreterOpts
}

// InterpreterOpts are options for creating a new instance.
type InterpreterOpts struct {
	// RunCommands is a string that is evaluated on startup.
	RunCommands string
}

// NewInterpreter creates a new interpreter.
func NewInterpreter(env *Environment, opts InterpreterOpts) (*Interpreter, error) {
	inst := Interpreter{
		env:  env,
		opts: opts,
	}

	inst.logger = log.New(inst.env.Terminal.Stderr, "", 0)

	inst.shParser = syntax.NewParser(
		syntax.KeepComments(false),
		syntax.Variant(syntax.LangBash), // we love bash!
	)

	shRunner, err := interp.New(
		interp.OpenHandler(func(ctx context.Context, path string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return env.Filesystem.OpenFile(path, flag, perm)
		}),
		interp.StatHandler(func(ctx context.Context, name string, followSymlinks bool) (fs.FileInfo, error) {
			return stdfs.Stat(env.Filesystem, name)
		}),
		interp.ReadDirHandler(func(ctx context.Context, path string) ([]fs.FileInfo, error) {
			entries, err := stdfs.ReadDir(env.Filesystem, path)
			if err != nil {
				return nil, err
			}

			infos := make([]fs.FileInfo, len(entries))
			for i, info := range entries {
				infos[i], err = info.Info()
				if err != nil {
					return nil, fmt.Errorf("%q: %v", info, err)
				}
			}

			return infos, nil
		}),
		interp.StdIO(inst.env.Terminal.Stdin, inst.env.Terminal.Stdout, inst.env.Terminal.Stderr),
		interp.ExecHandler(inst.execHandler),
		interp.Env(expand.ListEnviron(
			"HOME=/",
			"SITE=libdb.so",
		)),
		interp.RunnerOption(func(r *interp.Runner) error {
			r.Dir = "/"
			return nil
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init shell runner")
	}

	inst.shRunner = shRunner

	inst.prompter = liner.NewStateStdin(
		inst.env.Terminal.IO.Stdin,
		func() (row, col uint16, ok bool) {
			q := inst.env.Terminal.Query()
			return uint16(q.Width), uint16(q.Height), true
		},
	)

	return &inst, nil
}

// Close closes/stops the console instance.
func (inst *Interpreter) Close() error {
	var errs []error
	for _, fn := range inst.closes {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}
	return stderrors.Join(errs...)
}

// Terminal returns the terminal for the interpreter.
func (inst *Interpreter) Terminal() *Terminal {
	return inst.env.Terminal
}

// Run runs the console loop. It blocks until the context is canceled or until a
// fatal/unrecoverable error occurs.
func (inst *Interpreter) Run(ctx context.Context) error {
	ctx = context.WithValue(ctx, environmentKey, inst.env)
	ctx = context.WithValue(ctx, loggerKey, inst.logger)

	inst.exec(ctx, inst.opts.RunCommands)

	for {
		line, err := inst.prompter.Prompt(inst.prompt())
		if err != nil {
			return errors.Wrap(err, "failed to read line")
		}

		if line == "" {
			continue
		}

		inst.prompter.AppendHistory(line)
		inst.exec(ctx, line)
	}

	// interactiveFunc := func(stmts []*syntax.Stmt) bool {
	// 	if inst.shParser.Incomplete() {
	// 		fmt.Fprintf(inst.env.Terminal, "> ")
	// 		return true
	// 	}
	//
	// 	for _, stmt := range stmts {
	// 		if err := inst.shRunner.Run(ctx, stmt); err != nil {
	// 			inst.logger.Println(err)
	// 		}
	// 		if inst.shRunner.Exited() || ctx.Err() != nil {
	// 			return false
	// 		}
	// 	}
	//
	// 	fmt.Fprintf(inst.env.Terminal, inst.prompt())
	// 	return true
	// }
	//
	// fmt.Fprintf(inst.env.Terminal, inst.prompt())
	// return inst.shParser.Interactive(inst.env.Terminal.Stdin, interactiveFunc)
}

// Exec executes the given command string as if the user did it.
func (inst *Interpreter) Exec(ctx context.Context, line string) {
	ctx = context.WithValue(ctx, environmentKey, inst.env)
	ctx = context.WithValue(ctx, loggerKey, inst.logger)

	// By that, we just print the prompt and execute the command.
	inst.env.Print(inst.prompt(), line)
	inst.exec(ctx, line)
}

func (inst *Interpreter) exec(ctx context.Context, line string) {
	shFile, err := inst.shParser.Parse(strings.NewReader(line), "")
	if err != nil {
		inst.logger.Printf("error parsing: %v", err)
		return
	}

	if err := inst.shRunner.Run(ctx, shFile); err != nil {
		inst.logger.Printf("error running: %v", err)
		return
	}
}

func (inst *Interpreter) prompt() string {
	return "$ "
}

func (inst *Interpreter) execHandler(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	prog, ok := inst.env.Programs[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %q", args[0])
	}

	return prog.Run(ctx, inst.env, args)
}
