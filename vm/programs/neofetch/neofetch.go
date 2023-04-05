package neofetch

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	_ "embed"

	"github.com/fatih/color"
	"github.com/leaanthony/go-ansi-parser"
	"libdb.so/vm"
	"libdb.so/vm/programs"
)

var overrideRev string

// OverrideGitRevision overrides the git revision that is printed by neofetch.
func OverrideGitRevision(rev string) {
	overrideRev = rev
}

func init() {
	programs.Register(program{})
	color.NoColor = false
}

//go:embed me.sixel
var meSIXEL []byte

type program struct{}

func (program) Name() string { return "neofetch" }

func (program) Usage() string { return "print my information" }

func (program) Run(ctx context.Context, env vm.Environment, args []string) error {
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

	printLink(&b, spcolor("libdb.so", color.FgMagenta, color.Bold), "https://libdb.so/libdb.so")
	b.WriteByte('\n')
	b.WriteByte('\n')

	b.WriteByte('\n')

	fpcolor(&b, "Blog: ", color.FgHiMagenta, color.Bold)
	printLink(&b, "b.libdb.so", "https://b.libdb.so")
	b.WriteByte('\n')

	fpcolor(&b, "GitHub: ", color.FgHiCyan, color.Bold)
	printLink(&b, "diamondburned", "https://github.com/diamondburned")
	b.WriteByte('\n')

	fpcolor(&b, "Mastodon: ", color.FgHiBlue, color.Bold)
	printLink(&b, "@diamond@hachyderm.io", "https://hachyderm.io/@diamond")
	b.WriteByte('\n')

	b.WriteByte('\n')

	fmt.Fprintln(&b, spcolor("Go:", color.FgCyan), strings.Replace(runtime.Version(), "go", "v", 1))
	fmt.Fprintln(&b, spcolor("GOOS:", color.FgCyan), runtime.GOOS)
	fmt.Fprintln(&b, spcolor("GOARCH:", color.FgCyan), runtime.GOARCH)
	fmt.Fprintln(&b, spcolor("NumCPU:", color.FgCyan), cpuCount())
	if rev := programRev(); rev != "" {
		fmt.Fprintln(&b, spcolor("Version:", color.FgCyan), rev)
	}

	b.WriteByte('\n')

	printFgColors(&b, 0, 8)
	printFgColors(&b, 8, 16)
	return b.String()
}

func spcolor(s string, cs ...color.Attribute) string {
	return color.New(cs...).Sprint(s)
}

func fpcolor(w io.Writer, s string, cs ...color.Attribute) {
	color.New(cs...).Fprint(w, s)
}

func printFgColors(b *strings.Builder, from, to int) {
	for fg := from; fg <= to; fg++ {
		fmt.Fprintf(b, "\x1b[38;5;%dm%s\033[0m", fg, "███")
	}
	b.WriteByte('\n')
}

func cpuCount() string {
	ncpu := runtime.NumCPU()
	scpu := ""
	if ncpu == 1 {
		scpu = fmt.Sprintf("%d core", ncpu)
	} else {
		scpu = fmt.Sprintf("%d cores", ncpu)
	}

	ngoros := runtime.NumGoroutine()
	if ngoros == 1 {
		scpu += fmt.Sprintf(" (1 active goroutine)")
	} else {
		scpu += fmt.Sprintf(" (%d active goroutines)", ngoros)
	}

	return scpu
}

func programRev() string {
	vcs := "git"
	rev := overrideRev

	if rev == "" {
		build, ok := debug.ReadBuildInfo()
		if !ok {
			return ""
		}

		setting := func(k string) string {
			for _, setting := range build.Settings {
				if setting.Key == k {
					return setting.Value
				}
			}
			return ""
		}

		vcs = setting("vcs")
		rev = setting("vcs.revision")
		if vcs+rev == "" {
			return ""
		}
	}

	if len(rev) > 7 {
		rev = rev[:7]
	}

	return fmt.Sprintf("%s revision %s", vcs, rev)
}

var ansiLinkRe = regexp.MustCompile(`(?m)\x1b]8;;([^\x1b]*)\x1b\\([^\x1b]*)\x1b]8;;\x1b\\`)

const ansiLinkf = "\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\"

func init() {
	// keep in sync with the regex
	if !ansiLinkRe.MatchString(ansiLinkf) {
		panic("ansiLinkf is not in sync with ansiLinkRe")
	}
}

func printLink(b *strings.Builder, text, url string) {
	fmt.Fprintf(b, ansiLinkf, url, text)
}

func printInfo(env vm.Environment, str string) {
	const pad = 2
	const up = 16
	const right = 36

	lines := strings.Split(str, "\n")
	var maxLine int
	for _, line := range lines {
		// Replace all ANSI links with their text.
		line = ansiLinkRe.ReplaceAllString(line, "$2")
		// Count the length excluding ANSI codes.
		llen, _ := ansi.Length(line)
		if llen > maxLine {
			maxLine = llen
		}
	}

	if q := env.Terminal.Query(); q.Width > (right + pad + maxLine) {
		// go up
		env.Printf("\x1b[%dA", up)

		for _, line := range lines {
			env.Printf("\x1b[%dC", right)
			env.Println(line)
		}

		// go down delta lines
		if len(lines) < up {
			env.Printf("\x1b[%dB", up-len(lines))
		}
	} else {
		// No space, just print it.
		env.Println()
		env.Println(str)
	}
}
