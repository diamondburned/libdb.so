package console

import (
	"context"
	"fmt"
	"log"

	stderrors "errors"

	"github.com/pkg/errors"
	"libdb.so/console/fs"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Interpreter is the main interpreter for the larger program. It manages
// prompting the user, printing to console, and running programs.
type Interpreter struct {
	shParser *syntax.Parser
	shRunner *interp.Runner
	logger   *log.Logger
	term     *Terminal
	progs    map[string]Program
	closes   []func() error
}

// InterpreterOpts are options for creating a new instance.
type InterpreterOpts struct {
	// MakeRaw puts the terminal in raw mode. This is useful for running the
	// program in a sub-terminal.
	MakeRaw bool
	// Programs are the programs that are available to run.
	Programs []Program
	// Filesystem is the filesystem to use. If nil, then the OS filesystem is
	// used.
	Filesystem fs.FS // TODO
}

// NewInterpreter creates a new interpreter.
func NewInterpreter(io IO, opts InterpreterOpts) (*Interpreter, error) {
	var inst Interpreter

	inst.term = &Terminal{
		IO:    io,
		query: &terminalQueryUpdater{},
	}

	inst.progs = make(map[string]Program, len(opts.Programs))
	for _, prog := range opts.Programs {
		inst.progs[prog.Name()] = prog
	}

	inst.logger = log.New(io.Stderr, "console: ", 0)

	inst.shParser = syntax.NewParser(
		syntax.KeepComments(false),
		syntax.Variant(syntax.LangBash), // we love bash!
	)

	shRunner, err := interp.New(
		// TODO: interp.OpenHandler
		// TODO: interp.StatHandler
		// TODO: interp.ReadDirHandler
		interp.StdIO(io.Stdin, io.Stdout, io.Stderr),
		interp.ExecHandler(inst.execHandler),
		interp.Env(expand.ListEnviron(
			"HOME=/",
			"SITE=libdb.so",
		)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init shell runner")
	}

	inst.shRunner = shRunner

	if opts.MakeRaw {
		undo, err := io.makeRaw()
		if err != nil {
			return nil, err
		}

		inst.closes = append(inst.closes, undo)
	}

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

// UpdateTerminal updates the terminal with the given query.
func (inst *Interpreter) UpdateTerminal(query TerminalQuery) {
	inst.term.query.set(query)
}

// Run runs the console loop. It blocks until the context is canceled or until a
// fatal/unrecoverable error occurs.
func (inst *Interpreter) Run(ctx context.Context) error {
	interactiveFunc := func(stmts []*syntax.Stmt) bool {
		if inst.shParser.Incomplete() {
			fmt.Fprintf(inst.term, "> ")
			return true
		}

		for _, stmt := range stmts {
			if err := inst.shRunner.Run(ctx, stmt); err != nil {
				inst.logger.Println(err)
			}
			if inst.shRunner.Exited() || ctx.Err() != nil {
				return false
			}
		}

		fmt.Fprintf(inst.term, inst.prompt())
		return true
	}

	fmt.Fprintf(inst.term, inst.prompt())
	return inst.shParser.Interactive(inst.term.Stdin, interactiveFunc)
}

func (inst *Interpreter) prompt() string {
	return "$ "
}

func (inst *Interpreter) execHandler(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	prog, ok := inst.progs[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %q", args[0])
	}

	return prog.Run(inst.term, args)
}
