package global

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/lucasb-eyer/go-colorful"
	"gitlab.com/diamondburned/dotfiles/Scripts/lineprompt/lineprompt"
	"libdb.so/vm"
	"libdb.so/vm/rwfs"
	"libdb.so/vm/rwfs/kvfs"

	_ "libdb.so/vm/programs/coreutils"
	_ "libdb.so/vm/programs/hewwo"
	_ "libdb.so/vm/programs/neofetch"
	_ "libdb.so/vm/programs/nsfw"
	_ "libdb.so/vm/programs/resume"
	_ "libdb.so/vm/programs/sixel"
	_ "libdb.so/vm/programs/spew"
	_ "libdb.so/vm/programs/termio"
	_ "libdb.so/vm/programs/vars"
	_ "libdb.so/vm/programs/webring"
)

const InitialCwd = "/"

//go:embed shellrc
var shellrc []byte

// RootFS is the filesystem that contains default read-only files, such as the
// shellrc file.
var RootFS = rwfs.ReadOnlyFS(kvfs.New(kvfs.MemoryStorageFromExisting(
	map[string]kvfs.StoredValue{
		"/.shellrc": kvfs.StoredFile{Data: shellrc},
	},
)))

var InitialEnv = vm.EnvironFromMap(map[string]string{
	"TERM":  "xterm-256color",
	"HOME":  "/",
	"SITE":  "libdb.so",
	"SHELL": "github.com/mvdan/sh/v3",
	"SHLVL": "1",
})

// PromptMonochrome is a monochromic prompter.
func PromptMonochrome(env vm.Environment) string {
	return fmt.Sprintf("\n$ libdb.so @ %s\n―❤―▶ ", env.Cwd)
}

var transBlend = []colorful.Color{
	rgb(85, 205, 252),
	rgb(147, 194, 255),
	rgb(200, 181, 245),
	rgb(234, 171, 217),
	rgb(247, 148, 168),
}

// PromptColored is a colorful prompter.
func PromptColored() vm.PromptFunc {
	var s strings.Builder
	s.Grow(1024)

	var state struct {
		Column int
		Dir    string
	}

	return func(env vm.Environment) string {
		q := env.Terminal.Query()
		if q.Width == state.Column && state.Dir == env.Cwd {
			return s.String()
		}

		s.Reset()

		lineprompt.Blend(&s, "", q.Width, transBlend, lineprompt.Opts{
			LOD:       15,
			Underline: true,
		})
		s.WriteByte('\n')

		line1 := fmt.Sprintf("$ libdb.so @ %s", env.Cwd)
		if len(line1) > q.Width {
			line1 = line1[:q.Width]
		}

		lineprompt.Blend(&s, line1, q.Width, transBlend, lineprompt.Opts{
			LOD:       15,
			Underline: false,
		})
		s.WriteByte('\n')

		s.WriteString("\033[38;2;85;205;252m―❤―\033[0m\033[38;2;247;157;208m▶\033[m ")

		state.Column = q.Width
		state.Dir = env.Cwd
		return s.String()
	}
}

func rgb(r, g, b uint8) colorful.Color {
	return colorful.Color{
		R: float64(r) / 255,
		G: float64(g) / 255,
		B: float64(b) / 255,
	}
}
