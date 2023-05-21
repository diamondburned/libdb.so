package webring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	"libdb.so/libwebring-go"
	"libdb.so/vm"
	"libdb.so/vm/internal/ansi"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(app))
}

var webringURLs = []string{
	"https://raw.githubusercontent.com/diamondburned/acmfriends-webring/%3C3-spring-2023/webring.json",
}

var app = cli.App{
	Name:      "webring",
	Usage:     "print the webrings of this site",
	UsageText: "webring [name]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "print as JSON",
		},
		&cli.IntFlag{
			Name:    "width",
			Aliases: []string{"w"},
			Usage:   "width of the output",
		},
	},
	Action: func(c *cli.Context) error {
		return Print(c.Context, Opts{
			JSON:       c.Bool("json"),
			Width:      c.Int("width"),
			FilterName: c.Args().First(),
		})
	},
}

type Opts struct {
	JSON       bool
	Width      int
	FilterName string
}

// Print prints the webring.
func Print(ctx context.Context, opts Opts) error {
	env := vm.EnvironmentFromContext(ctx)
	log := vm.LoggerFromContext(ctx)

	webrings := make([]*libwebring.Data, 0, len(webringURLs))
	for _, url := range webringURLs {
		w, err := libwebring.FetchData(ctx, url)
		if err != nil {
			log.Println("cannot fetch webring:", err)
			continue
		}

		s, err := libwebring.FetchStatusForWebring(ctx, url)
		if err != nil {
			log.Println("cannot check webring:", err)
			continue
		}

		w.Ring = w.Ring.ExcludeAnomalies(s.Anomalies)
		webrings = append(webrings, w)
	}

	if len(webrings) == 0 && len(webringURLs) > 0 {
		return errors.New("no webrings found")
	}

	if opts.FilterName != "" {
		webrings = []*libwebring.Data{findWebring(webrings, opts.FilterName)}
		if webrings[0] == nil {
			return fmt.Errorf("webring %q not found", opts.FilterName)
		}
	}

	if opts.JSON {
		e := json.NewEncoder(env.Terminal.Stdout)
		e.SetIndent("", "  ")
		return e.Encode(webrings)
	}

	width := opts.Width
	if width == 0 {
		width = env.Terminal.Query().Width
	}

	ringPrefix := "part of the "
	if width < 50 {
		ringPrefix = ""
	}

	for i, w := range webrings {
		ring := w.Name
		if w.Root != "" {
			ring = ansi.Link(w.Name, w.Root)
		}

		mine := findLink(w)
		if mine == nil {
			log.Println("diamond is no longer in webring", w.Name)
			continue
		}

		left, right := w.Ring.SurroundingLinks(*mine)
		printAligned(env, width,
			"← "+ansi.Link(left.Name, ensureScheme(left.Link)),
			"─── "+ringPrefix+ring+" webring ───",
			ansi.Link(right.Name, ensureScheme(right.Link))+" →",
		)

		if i != len(webrings)-1 {
			env.Println()
		}
	}

	return nil
}

func ensureScheme(url string) string {
	if !strings.Contains(url, "://") {
		return "https://" + url
	}
	return url
}

func findLink(webring *libwebring.Data) *libwebring.Link {
	for i, link := range webring.Ring {
		if link.Name == "diamond" || link.Link == "libdb.so" {
			return &webring.Ring[i]
		}
	}
	return nil
}

func findWebring(webrings []*libwebring.Data, name string) *libwebring.Data {
	for _, w := range webrings {
		if w.Name == name {
			return w
		}
	}
	return nil
}

func printAligned(env vm.Environment, width int, left, center, right string) {
	w := env.Terminal.Stdout

	leftWidth := ansi.StringWidth(left)
	rightWidth := ansi.StringWidth(right)
	centerWidth := ansi.StringWidth(center)
	halfCenterWidth := centerWidth / 2

	totalWidth := ansi.StringWidth(left) + ansi.StringWidth(center) + ansi.StringWidth(right) + 3
	if totalWidth > width {
		centerPad := max((width-centerWidth)/2, 0)
		rightPad := max(width-rightWidth, 0)

		fmt.Fprintln(w, left)

		fmt.Fprint(w, strings.Repeat(" ", centerPad))
		fmt.Fprintln(w, center)

		fmt.Fprint(w, strings.Repeat(" ", rightPad))
		fmt.Fprintln(w, right)
		return
	}

	pivot := width / 2
	leftPad := max(pivot-halfCenterWidth-leftWidth, 1)
	rightPad := max(width-leftWidth-leftPad-centerWidth-rightWidth, 1)

	fmt.Fprint(w, left)
	fmt.Fprint(w, strings.Repeat(" ", leftPad))
	fmt.Fprint(w, center)
	fmt.Fprint(w, strings.Repeat(" ", rightPad))
	fmt.Fprint(w, right)

	fmt.Fprintln(w)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
