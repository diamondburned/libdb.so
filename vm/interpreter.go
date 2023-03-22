package vm

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"sort"
	"strings"
	"text/tabwriter"

	stderrors "errors"

	"github.com/pkg/errors"
	"libdb.so/vm/internal/liner"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// PromptFunc is a function that returns the prompt string.
type PromptFunc func(Environment) string

// Interpreter is the main interpreter for the larger program. It manages
// prompting the user, printing to console, and running programs.
type Interpreter struct {
	shParser    *syntax.Parser
	shRunner    *interp.Runner
	shExpandCfg *expand.Config
	prompter    *liner.State
	logger      *log.Logger
	env         *Environment
	closes      []func() error
	opts        InterpreterOpts
	progNames   []string
}

// InterpreterOpts are options for creating a new instance.
type InterpreterOpts struct {
	// RunCommands is a string that is evaluated on startup.
	RunCommands string
	// Prompt is a function that returns the prompt string.
	Prompt PromptFunc
	// IgnoreEOF, if true, will ignore EOF errors and continue prompting as
	// usual.
	IgnoreEOF bool
}

var builtinCommands = []string{
	"true", "false", "exit", "set", "shift", "unset", "echo", "printf", "pwd",
	"cd", "source", "command", "umask", "alias", "unalias", "eval", "test",
	"exec", "read", "readarray", "shopt",
}

// NewInterpreter creates a new interpreter.
func NewInterpreter(env *Environment, opts InterpreterOpts) (*Interpreter, error) {
	inst := Interpreter{
		env:  env,
		opts: opts,
	}

	inst.progNames = make([]string, 0, len(env.Programs))
	inst.progNames = append(inst.progNames, builtinCommands...)
	for name := range env.Programs {
		inst.progNames = append(inst.progNames, name)
	}
	sort.Strings(inst.progNames)

	if inst.opts.Prompt == nil {
		inst.opts.Prompt = func(Environment) string { return "$ " }
	}

	inst.logger = log.New(inst.env.Terminal.Stderr, "", 0)

	readDir := func(path string) ([]fs.FileInfo, error) {
		entries, err := fs.ReadDir(env.Filesystem, path)
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
	}

	inst.shExpandCfg = &expand.Config{
		Env:     inst.env.Environ,
		ReadDir: readDir,
		// TODO: CmdSubst
	}

	inst.shParser = syntax.NewParser(
		syntax.KeepComments(false),
		syntax.Variant(syntax.LangBash), // we love bash!
	)

	shRunner, err := interp.New(
		// TODO: CmdSubst
		interp.OpenHandler(func(ctx context.Context, path string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			inst.updateEnv()
			return env.Filesystem.OpenFile(env.JoinCwd(path), flag, perm)
		}),
		interp.StatHandler(func(ctx context.Context, name string, followSymlinks bool) (fs.FileInfo, error) {
			inst.updateEnv()
			return fs.Stat(env.Filesystem, name)
		}),
		interp.ReadDirHandler(func(ctx context.Context, path string) ([]fs.FileInfo, error) {
			inst.updateEnv()
			return readDir(path)
		}),
		interp.StdIO(inst.env.Terminal.Stdin, inst.env.Terminal.Stdout, inst.env.Terminal.Stderr),
		interp.ExecHandler(inst.execHandler),
		interp.CallHandler(inst.callHandler),
		interp.Env(inst.env.Environ),
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
			return uint16(q.Height), uint16(q.Width), true
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
func (inst *Interpreter) Terminal() Terminal {
	return inst.env.Terminal
}

// Run runs the console loop. It blocks until the context is canceled or until a
// fatal/unrecoverable error occurs.
func (inst *Interpreter) Run(ctx context.Context) error {
	ctx = context.WithValue(ctx, environmentKey, inst.env)
	ctx = context.WithValue(ctx, loggerKey, inst.logger)

	inst.exec(ctx, inst.opts.RunCommands)
	inst.prompter.SetWordCompleter(inst.wordCompleter(ctx))
	inst.prompter.SetTabCompletionStyle(liner.TabPrints)
	// Fully aborting lets us draw the entire prompt instead of just the
	// incomplete line.
	inst.prompter.SetCtrlCAborts(true)

	for {
		inst.env.Cwd = inst.shRunner.Dir
		inst.env.Environ = inst.shRunner.Env
		prompt := inst.opts.Prompt(*inst.env)

		// Support multiline prompts by splitting on newlines and printing each
		// line separately except the last one.
		promptLines := strings.SplitAfter(prompt, "\n")
		for _, line := range promptLines[:len(promptLines)-1] {
			inst.env.Print(line)
		}

		line, err := inst.prompter.Prompt(promptLines[len(promptLines)-1])
		if err != nil {
			switch {
			case errors.Is(err, liner.ErrPromptAborted):
				// Ctrl+C is pressed; redraw entire prompt.
				continue
			case errors.Is(err, io.EOF):
				if inst.opts.IgnoreEOF {
					inst.env.Println()
					continue
				}
				return nil // break with no error
			default:
				return errors.Wrap(err, "failed to read line")
			}
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

func (inst *Interpreter) callHandler(ctx context.Context, args []string) ([]string, error) {
	inst.updateEnv()
	return args, nil
}

func (inst *Interpreter) execHandler(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return nil
	}

	handler := interp.HandlerCtx(ctx)
	inst.updateEnv()

	env := *inst.env // this is ephemeral
	env.Cwd = handler.Dir
	env.Terminal = env.Terminal.WithIO(IO{
		Stdin:  io.NopCloser(handler.Stdin),
		Stdout: handler.Stdout,
		Stderr: handler.Stderr,
	})

	switch args[0] {
	case "help":
		return inst.help(env)
	}

	prog, ok := inst.env.Programs[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %q", args[0])
	}

	ctx = context.WithValue(ctx, environmentKey, &env)
	return prog.Run(ctx, env, args)
}

func (inst *Interpreter) wordCompleter(ctx context.Context) func(string, int) (string, []string, string) {
	return func(line string, pos int) (head string, completions []string, tail string) {
		shf, err := inst.shParser.Parse(strings.NewReader(line), "")
		if err != nil {
			// cannot be parsed, ignore
			return
		}

		// We keep track of the call that the cursor is on as well as the shell
		// word within that statement that we're trying to autocomplete.
		//
		// When we DFS a line like "echo && hi a", we'll get something like this:
		//
		//  - Stmt
		//    - BinaryCmd
		//      - Stmt
		//        - CallExpr
		//          - Word
		//            - Lit
		//      - Stmt
		//        - CallExpr
		//          - Word
		//            - Lit
		//          - Word
		//            - Lit
		//
		// Our cursor might be sandwiched between some of these nodes, so we can
		// check for their positions to extract them.
		var cursorExpr *syntax.CallExpr
		var cursorWord *syntax.Word

		syntax.Walk(shf, func(n syntax.Node) (deeper bool) {
			if n != nil {
				// log.Printf("walking %T (%d, %d)", n, n.Pos().Offset(), n.End().Offset())
			}

			switch n := n.(type) {
			case nil:
				// we're traversing up the tree
				return false
			case *syntax.CallExpr:
				// We're checking if the cursor is inside this call expression,
				// but we also check one more case where the rest of the string
				// is just whitespaces. This suggests that the cursor is still
				// in this call expression.
				if shNodeCovers(n, pos) || justSpace(line[n.End().Offset():]) {
					cursorExpr = n
				}
			case *syntax.Word:
				if shNodeCovers(n, pos) {
					cursorWord = n
				}
			}

			return true
		})

		if cursorExpr == nil {
			if line == "" {
				// empty line, so we just autocomplete for programs.
				completions = inst.programAutocomplete("")
			}
			return
		}

		var cursorValue string
		var ix int

		if cursorWord == nil {
			// Handle a particular case where the cursor is at the end and
			// happens to be on a space. In that case, we want to list all
			// possible words.
			if pos == len(line) && line[pos-1] == ' ' {
				ix = len(cursorExpr.Args)
				head = line[:pos]
				tail = ""
			} else {
				log.Println("cursor is in a word without a cursorWord")
				return
			}
		} else {
			cursorValue, err = expand.Literal(inst.shExpandCfg, cursorWord)
			if err != nil {
				return
			}

			ix = wordIndexWithin(cursorExpr, cursorWord)
			if ix == -1 {
				log.Println("vm: interpreter: possible bug: word not found in call expression")
				return
			}

			head = line[:cursorWord.Pos().Offset()]
			tail = line[cursorWord.End().Offset():]
		}

		if ix == 0 {
			// arg0, we'll do program autocompletion.
			completions = inst.programAutocomplete(cursorValue)
			return
		}

		// Not arg0, so we'll need to know what arg0 actually is in order to check
		// if it is a program that we can autocomplete.
		firstValue, err := expand.Literal(inst.shExpandCfg, cursorExpr.Args[0])
		if err != nil {
			return
		}

		prog, found := inst.env.Programs[firstValue]
		autocompleter, isAutocompleter := prog.(ProgramAutocompleter)
		if found && isAutocompleter {
			// We should guarantee that we can expand words up to the
			// cursor. Don't bother with the rest if we have no choice.
			words := make([]string, 0, len(cursorExpr.Args))
			for i, w := range cursorExpr.Args {
				w, err := expand.Literal(inst.shExpandCfg, w)
				if err != nil {
					if i < ix {
						// we couldn't decode up to the cursor, ignore
						return
					}
					// we can't decode the rest, but we can still
					// autocomplete up to the cursor
					break
				}
				words = append(words, w)
			}

			// Is this even needed? Is this overkill? I think it is, but
			// whatever.
			env := Environment{
				Terminal:   inst.env.Terminal.WithIO(NoIOExceptStderr(inst.env.Terminal.Stderr)),
				Filesystem: inst.env.Filesystem,
				Cwd:        inst.shRunner.Dir,
				Programs:   inst.env.Programs,
				Environ:    inst.env.Environ,
			}

			completions = autocompleter.Autocomplete(ctx, env, words, ix)
			return
		}

		// fallback: do file autocompletion
		files, err := fs.ReadDir(inst.env.Filesystem, inst.env.Cwd)
		if err != nil {
			// We can't read PWD. That's pretty bad, but we can't do much.
			return
		}

		for _, f := range files {
			// If we're autocompleting cd, then we only want directories.
			if firstValue == "cd" && !f.IsDir() {
				continue
			}
			if strings.HasPrefix(f.Name(), cursorValue) {
				completions = append(completions, f.Name())
			}
		}

		return
	}
}

func (inst *Interpreter) updateEnv() {
	inst.env.Cwd = inst.shRunner.Dir
	inst.env.Environ = inst.shRunner.Env
}

func (inst *Interpreter) programAutocomplete(word string) []string {
	var completions []string
	for _, name := range inst.progNames {
		if strings.HasPrefix(name, word) {
			completions = append(completions, name)
		}
	}
	return completions
}

// wordIndexWithin returns the index of the word within the statement. If word
// is not found, -1 is returned.
func wordIndexWithin(call *syntax.CallExpr, word *syntax.Word) int {
	for i, w := range call.Args {
		if w == word {
			return i
		}
	}
	return -1
}

func shNodeCovers(node syntax.Node, pos int) bool {
	return node.Pos().Offset() <= uint(pos) && uint(pos) <= node.End().Offset()
}

func justSpace(str string) bool {
	return strings.TrimSpace(str) == ""
}

func (inst *Interpreter) help(env Environment) error {
	fmt.Fprint(env.Terminal.Stdout, "Available commands:\n\n")
	w := tabwriter.NewWriter(env.Terminal.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "help\tShow this help\n")
	for _, name := range inst.progNames {
		prog, ok := inst.env.Programs[name]
		if !ok {
			continue
		}
		if usager, ok := prog.(ProgramUsager); ok {
			fmt.Fprintf(w, "%s\t%s\n", usager.Name(), usager.Usage())
		} else {
			fmt.Fprintf(w, "%s\t\n", prog.Name())
		}
	}
	return w.Flush()
}
