package coreutils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"path"
	"strings"
	"text/tabwriter"
	"time"

	stderrors "errors"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/ansi"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/internal/vmutil"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(ls))
}

var ls = cli.App{
	Name:      "ls",
	Usage:     "list directory contents",
	UsageText: `ls [OPTION]... [FILE]...`,
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
		args := c.Args().Slice()
		if len(args) == 0 {
			args = []string{"."}
		}

		var errs []error
		for _, arg := range args {
			if err := ls_(c, arg, len(args) > 1); err != nil {
				errs = append(errs, err)
			}
		}

		return stderrors.Join(errs...)
	},
}

func ls_(c *cli.Context, arg string, multiple bool) error {
	env := vm.EnvironmentFromContext(c.Context)
	path := path.Join(env.Cwd, arg)

	stat, err := fs.Stat(env.Filesystem, path)
	if err != nil {
		return errors.Wrap(err, "stat")
	}

	var ents []fs.DirEntry

	if stat.IsDir() {
		ents, err = fs.ReadDir(env.Filesystem, path)
		if err != nil {
			return errors.Wrap(err, "readdir")
		}
	} else {
		ents = []fs.DirEntry{
			fakeDirEntry{stat},
		}
	}

	if multiple {
		fmt.Fprintln(env.Terminal.Stdout, arg+":")
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
			var modTime time.Time
			var mode fs.FileMode
			var size int64

			s, err := ent.Info()
			if err != nil {
				log.Println("stat:", err)
			} else {
				modTime = s.ModTime()
				mode = s.Mode()
				size = s.Size()
			}

			fmt.Fprintf(w,
				"%s\t%d\t%s\t%s\n",
				printPerm(mode), size, printTime(modTime), printName(env, ent),
			)
		}

		return w.Flush()
	}

	for _, ent := range ents {
		fmt.Fprintln(c.App.Writer, printName(env, ent))
	}

	return nil
}

func printName(env vm.Environment, dirEntry fs.DirEntry) string {
	name := dirEntry.Name()
	if env.HasTerminal {
		name = ansi.Link(name, vmutil.MakeTerminalWriteURI(name))
		if dirEntry.IsDir() {
			return color.New(color.FgBlue, color.Bold).Sprint(name)
		}
	}
	return name
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
