package coreutils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(ls))
}

var ls = cli.App{
	Name:  "ls",
	Usage: "list directory contents",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "do not ignore entries starting with .",
		},
		&cli.BoolFlag{
			Name:    "long",
			Aliases: []string{"l"},
			Usage:   "use a long listing format",
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "use a JSON listing format",
		},
	},
	Action: func(c *cli.Context) error {
		env := vm.EnvironmentFromContext(c.Context)

		arg := c.Args().First()
		if arg == "" {
			arg = env.Cwd
		}

		stat, err := fs.Stat(env.Filesystem, arg)
		if err != nil {
			return errors.Wrap(err, "stat")
		}

		var ents []fs.DirEntry

		if stat.IsDir() {
			ents, err = fs.ReadDir(env.Filesystem, arg)
			if err != nil {
				return errors.Wrap(err, "failed to read directory")
			}
		} else {
			ents = []fs.DirEntry{
				fakeDirEntry{stat},
			}
		}

		if !c.Bool("all") {
			filtered := ents[:0]
			for _, ent := range ents {
				if strings.HasPrefix(ent.Name(), ".") {
					continue
				}
				filtered = append(filtered, ent)
			}
			ents = filtered
		}

		if c.Bool("json") {
			lsEnts := make([]lsEntry, len(ents))
			for i, ent := range ents {
				lsEnts[i] = lsEntry{
					Name:  ent.Name(),
					Type:  ent.Type(),
					IsDir: ent.IsDir(),
				}
			}

			enc := json.NewEncoder(c.App.Writer)
			enc.SetIndent("", "  ")
			return enc.Encode(lsEnts)
		}

		if c.Bool("long") {
			w := tabwriter.NewWriter(c.App.Writer, 0, 0, 1, ' ', 0)

			for _, ent := range ents {
				var t time.Time
				var m fs.FileMode

				s, err := ent.Info()
				if err == nil {
					t = s.ModTime()
					m = s.Mode()
				}

				fmt.Fprintf(w,
					"%s\t%s\t%s\n",
					printPerm(m), printTime(t), printName(ent),
				)
			}

			return w.Flush()
		}

		for _, ent := range ents {
			fmt.Fprintln(c.App.Writer, printName(ent))
		}

		return nil
	},
}

func printName(dirEntry fs.DirEntry) string {
	if dirEntry.IsDir() {
		return color.New(color.FgBlue, color.Bold).Sprint(dirEntry.Name())
	}
	return dirEntry.Name()
}

func printPerm(mode fs.FileMode) string {
	var str [9]byte
	perm := mode.Perm()
	for i := 0; i < 9; i++ {
		if perm&(1<<uint(8-i)) != 0 {
			str[i] = "rwx"[i%3]
		} else {
			str[i] = '-'
		}
	}
	return string(str[:])
}

func printTime(time time.Time) string {
	return time.Format("Jan 02 15:04")
}

type fakeDirEntry struct {
	fs.FileInfo
}

var _ fs.DirEntry = fakeDirEntry{}

func (e fakeDirEntry) Name() string               { return e.FileInfo.Name() }
func (e fakeDirEntry) Type() fs.FileMode          { return e.FileInfo.Mode().Type() }
func (e fakeDirEntry) IsDir() bool                { return e.FileInfo.IsDir() }
func (e fakeDirEntry) Info() (fs.FileInfo, error) { return e.FileInfo, nil }

type lsEntry struct {
	Name  string      `json:"name"`
	Type  fs.FileMode `json:"type"`
	IsDir bool        `json:"is_dir"`
}
