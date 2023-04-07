package neofetch

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"text/tabwriter"

	_ "embed"

	"github.com/fatih/color"
	"github.com/leaanthony/go-ansi-parser"
	"github.com/lucasb-eyer/go-colorful"
	"gitlab.com/diamondburned/dotfiles/Scripts/lineprompt/lineprompt"
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

func colorfulHex(hex string) colorful.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(err)
	}
	return c
}

var transColors = []colorful.Color{
	colorfulHex("#55CDFC"),
	colorfulHex("#F7A8B8"),
	colorfulHex("#FFFFFF"),
}

func transBand(inverted bool) string {
	var s strings.Builder
	b := func(c colorful.Color) {
		const b = "█"
		s.Write(lineprompt.RGB(c.RGB255()))
		s.WriteString(b)
		s.WriteString(b)
	}

	if inverted {
		for i := len(transColors) - 1; i >= 0; i-- {
			b(transColors[i])
		}
	} else {
		for i := range transColors {
			b(transColors[i])
		}
	}

	s.Write(lineprompt.Reset)
	return s.String()
}

func info() string {
	var b strings.Builder
	b.WriteByte('\n')

	fmt.Fprintln(&b,
		transBand(false),
		spcolor("diamondburned", color.FgHiMagenta, color.Bold),
		spcolor("(she/they/it)", color.FgHiMagenta),
		transBand(true),
	)

	b.WriteByte('\n')

	b.WriteByte('\n')

	columnate2(&b,
		spcolor("Blog:", color.FgHiGreen, color.Bold),
		splink("b.libdb.so", "https://b.libdb.so"),

		spcolor("GitHub:", color.FgHiCyan, color.Bold),
		splink("diamondburned", "https://github.com/diamondburned"),

		spcolor("Mastodon:", color.FgHiBlue, color.Bold),
		splink("@diamond@hachyderm.io", "https://hachyderm.io/@diamond"),
	)

	b.WriteByte('\n')

	source := splink("GitHub", "https://github.com/diamondburned/libdb.so")
	if rev := programRev(); rev != "" {
		source += " (" + rev + ")"
	}

	columnate2(&b,
		spcolor("Go:", color.FgCyan),
		fmt.Sprintf("%s on %s/%s",
			strings.Replace(runtime.Version(), "go", "v", 1),
			runtime.GOOS, runtime.GOARCH),

		spcolor("NumCPU:", color.FgCyan),
		cpuCount(),

		spcolor("Source:", color.FgCyan),
		source,
	)

	b.WriteByte('\n')

	printFgColors(&b, 0, 8)
	printFgColors(&b, 8, 16)
	return b.String()
}

func columnate2(w io.Writer, values ...string) {
	columnate(w, 2, values...)
}

func columnate(w io.Writer, n int, values ...string) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	for i := 0; i < len(values); i += 2 {
		fmt.Fprintf(tw, "%s\t%s\n", values[i], values[i+1])
	}
	tw.Flush()
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

func splink(text, url string) string {
	return fmt.Sprintf(ansiLinkf, url, text)
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
