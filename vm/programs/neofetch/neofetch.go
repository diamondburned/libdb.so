package neofetch

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	_ "embed"

	"github.com/fatih/color"
	"libdb.so/vm"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(program{})
}

//go:embed me.sixel
var meSIXEL []byte

type program struct{}

func (program) Name() string { return "neofetch" }

func (program) Run(ctx context.Context, env *vm.Environment, args []string) error {
	if len(args) != 1 {
		return &vm.UsageError{Usage: "neofetch"}
	}

	env.Terminal.Write(meSIXEL)
	env.Terminal.Write([]byte("\n"))
	printInfo(env, info())
	return nil
}

func info() string {
	var b strings.Builder

	color.New(color.FgMagenta, color.Bold).Fprintf(&b, "libdb.so\n")
	color.New(color.FgMagenta, color.Bold).Fprintf(&b, "--------\n")
	b.WriteByte('\n')

	b.WriteByte('\n')

	b.WriteString("GitHub: ")
	printLink(&b, "diamondburned", "https://github.com/diamondburned")
	b.WriteByte('\n')

	b.WriteString("Mastodon: ")
	printLink(&b, "@diamond@hachyderm.io", "https://hachyderm.io/@diamond")
	b.WriteByte('\n')

	b.WriteByte('\n')

	build, _ := debug.ReadBuildInfo()
	fmt.Fprintln(&b, "Go version", build.GoVersion)
	fmt.Fprintln(&b, "GOOS:", runtime.GOOS)
	fmt.Fprintln(&b, "GOARCH:", runtime.GOARCH)
	fmt.Fprintln(&b, "NumCPU:", runtime.NumCPU())

	b.WriteByte('\n')

	printFgColors(&b, 0, 8)
	printFgColors(&b, 8, 16)
	return b.String()
}

func printFgColors(b *strings.Builder, from, to int) {
	for fg := from; fg <= to; fg++ {
		fmt.Fprintf(b, "\x1b[38;5;%dm%s\033[0m", fg, "███")
	}
	b.WriteByte('\n')
}

func printLink(b *strings.Builder, text, url string) {
	fmt.Fprintf(b, "\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

func printInfo(env *vm.Environment, str string) {
	const up = 16
	const right = 36

	// go up
	env.Printf("\x1b[%dA", up)

	lines := strings.Split(str, "\n")
	for _, line := range lines {
		env.Printf("\x1b[%dC", right)
		env.Println(line)
	}

	// go down delta lines
	if len(lines) < up {
		env.Printf("\x1b[%dB", up-len(lines))
	}
}
