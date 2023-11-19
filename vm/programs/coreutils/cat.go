package coreutils

import (
	"bufio"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/internal/nsfw"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(cat))
}

var cat = cli.App{
	Name:      "cat",
	Usage:     "concatenate files and print on the standard output",
	UsageText: `cat [FILE]...`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-sixel",
			Usage: "print sixel graphics as raw bytes (always true if no img2sixel)",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			if nsfw.IsEnabled() {
				env := vm.EnvironmentFromContext(c.Context)
				env.Println("nyaa~")
				return nil
			}
			return errors.New("no files given")
		}

		var failed bool
		for _, arg := range c.Args().Slice() {
			if !printFile(c, arg) {
				failed = true
			}
		}

		if failed {
			return errors.New("failed to print one or more files")
		}

		return nil
	},
}

func printFile(c *cli.Context, path string) bool {
	env := vm.EnvironmentFromContext(c.Context)
	log := vm.LoggerFromContext(c.Context)

	f, err := env.Open(path)
	if err != nil {
		log.Println("open:", err)
		return false
	}
	defer f.Close()

	mime, r, err := readMIME(f)
	if err != nil {
		log.Println("readMIME:", err)
		return false
	}

	switch mime {
	case "image/jpeg", "image/png":
		if c.Bool("no-sixel") {
			break
		}

		err := env.Execute(c.Context, env, "img2sixel", path)
		if vm.ErrorIsUnknownProgram(err) {
			// Just print as raw if img2sixel is not available.
			break
		}

		if err != nil {
			log.Println(err)
			return false
		}

		return true
	}

	if _, err = io.Copy(env.Terminal.Stdout, r); err != nil {
		log.Println("io.Copy:", err)
	}

	return true
}

func readMIME(r io.Reader) (string, io.Reader, error) {
	buf := bufio.NewReaderSize(r, 512)

	head, err := buf.Peek(512)
	if err != nil && err != io.EOF {
		return "", nil, errors.Wrap(err, "peek")
	}

	return http.DetectContentType(head), buf, nil
}
