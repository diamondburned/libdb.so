package vm

import (
	"context"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"libdb.so/vm/fs"
	"mvdan.cc/sh/v3/expand"
)

// Program defines a userspace program within our larger program.
type Program interface {
	// Name returns the name of the program.
	Name() string
	// Run runs the program. It should return an error if the program failed to
	// run. args is the arguments passed to the program; args[0] is the name of
	// the program.
	Run(ctx context.Context, env Environment, args []string) error
}

// ProgramAutocompleter defines a program that supports autocompletion.
type ProgramAutocompleter interface {
	Program
	// Autocomplete returns a list of autocompletions for the given arguments.
	// The function is given the full list of arguments, as well as the position
	// of the cursor relative to that list.
	Autocomplete(ctx context.Context, env Environment, args []string, cursor int) []string
}

// UsageError is an error that indicates the user provided invalid arguments.
type UsageError struct {
	Err   error // optional
	Usage string
}

func (err *UsageError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("invalid usage: %s (see: %s)", err.Err, err.Usage)
	}
	return fmt.Sprintf("usage: %s", err.Usage)
}

func (err *UsageError) Unwrap() error { return err.Err }

// EnvironFromMap returns an expand.Environ from a map.
func EnvironFromMap(envs map[string]string) expand.Environ {
	var environ []string
	for k, v := range envs {
		if strings.Contains(k, "=") {
			log.Panicf("invalid environment variable name %q", k)
		}
		environ = append(environ, fmt.Sprintf("%s=%s", k, v))
	}
	return expand.ListEnviron(environ...)
}

// Environment contains the environment for a program.
type Environment struct {
	// Terminal is the terminal to use.
	Terminal Terminal
	// Filesystem is the filesystem to use.
	Filesystem fs.FS
	// Cwd is the current working directory.
	Cwd string
	// Programs are the programs that are available to run.
	Programs map[string]Program
	// Environ is the environment variables.
	Environ expand.Environ
}

// Env returns the environment variable with the given key.
func (env *Environment) Env(key string) string {
	_, v := env.Environ.Get(key).Resolve(env.Environ)
	return v.Str
}

// JoinCwd joins the given path components and prepends the current working
// directory.
func (env *Environment) JoinCwd(paths ...string) string {
	return path.Join(append([]string{env.Cwd}, paths...)...)
}

// ListPrograms returns a list of all programs in alphabetical order.
func (env *Environment) ListPrograms() []string {
	progs := make([]string, 0, len(env.Programs))
	for name := range env.Programs {
		progs = append(progs, name)
	}
	sort.Slice(progs, func(i, j int) bool {
		return progs[i] < progs[j]
	})
	return progs
}

// Println prints a line to the terminal.
func (env *Environment) Println(v ...any) {
	fmt.Fprintln(env.Terminal.Stdout, v...)
}

// Printf prints a formatted line to the terminal.
func (env *Environment) Printf(f string, v ...any) {
	fmt.Fprintf(env.Terminal.Stdout, f, v...)
}

// Print prints to the terminal.
func (env *Environment) Print(v ...any) {
	fmt.Fprint(env.Terminal.Stdout, v...)
}