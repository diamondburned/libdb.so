package neofetch

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strings"
	"text/tabwriter"

	_ "embed"

	"github.com/fatih/color"
	"github.com/lucasb-eyer/go-colorful"
	"gitlab.com/diamondburned/dotfiles/Scripts/lineprompt/lineprompt"
	"libdb.so/vm"
	"libdb.so/vm/internal/ansi"
	"libdb.so/vm/internal/nsfw"
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

	pronouns := "(she/they)"
	if nsfw.IsEnabled() {
		pronouns = "(she/they/it)"
	}

	fmt.Fprintln(&b,
		transBand(false),
		spcolor("diamondburned", color.FgHiMagenta, color.Bold),
		spcolor(pronouns, color.FgHiMagenta),
		transBand(true),
	)

	b.WriteByte('\n')

	var mastodonCode string
	if nsfw.IsEnabled() {
		mastodonCode = ansi.Link(
			"@diamond@girlcock.club",
			"https://girlcock.club/@diamond",
		)
		mastodonCode += " (nsfw)"
	} else {
		mastodonCode = ansi.Link(
			"@diamond@tech.lgbt",
			"https://tech.lgbt/@diamond",
		)
	}

	source := ansi.Link("GitHub", "https://github.com/diamondburned/libdb.so")
	if rev := programRev(); rev != "" {
		source += " (" + rev + ")"
	}

	columnate2(&b,
		spcolor("Blog", color.FgHiYellow, color.Bold),
		ansi.Link("b.libdb.so", "https://b.libdb.so"),

		spcolor("Email", color.FgHiMagenta, color.Bold),
		ansi.Link("x@libdb.so", "mailto:x@libdb.so"),

		spcolor("GitHub", color.FgHiCyan, color.Bold),
		ansi.Link("diamondburned", "https://github.com/diamondburned"),

		spcolor("Matrix", color.FgHiRed, color.Bold),
		ansi.Link("@diamondburned:matrix.org", "https://matrix.to/#/@diamondburned:matrix.org"),

		spcolor("Discord", color.FgHiBlue, color.Bold),
		"@diamondburned",

		spcolor("Mastodon", color.FgHiGreen, color.Bold),
		mastodonCode,

		"", "",

		spcolor("Go", color.Reset, color.FgCyan),
		fmt.Sprintf("%s on %s/%s",
			strings.Replace(runtime.Version(), "go", "v", 1),
			runtime.GOOS, runtime.GOARCH),

		spcolor("Source", color.Reset, color.FgCyan),
		source,
	)

	b.WriteByte('\n')

	printFgColors(&b, 0, 8)
	printFgColors(&b, 8, 16)
	return b.String()
}

func columnate2(w io.Writer, values ...string) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
	for i := 0; i < len(values); i += 2 {
		if values[i] == "" && values[i+1] == "" {
			fmt.Fprintf(tw, "\t\n")
		} else {
			fmt.Fprintf(tw, "%s\t%s\n", values[i], values[i+1])
		}
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

func printInfo(env vm.Environment, str string) {
	const pad = 2
	const up = 16
	const right = 39

	lines := strings.Split(str, "\n")
	var maxLine int
	for _, line := range lines {
		if llen := ansi.StringWidth(line); llen > maxLine {
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
