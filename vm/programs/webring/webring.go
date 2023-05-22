package webring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

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
		if err == nil {
			w.Ring = w.Ring.ExcludeAnomalies(s.Anomalies)
		}

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
		ansi.PrintAligned(env, width,
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
