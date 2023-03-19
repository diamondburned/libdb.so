package global

import (
	"bytes"
	"fmt"

	_ "embed"

	"github.com/lucasb-eyer/go-colorful"
	"gitlab.com/diamondburned/dotfiles/Scripts/lineprompt/lineprompt"
	"libdb.so/vm"

	_ "libdb.so/vm/programs/coreutils"
	_ "libdb.so/vm/programs/hewwo"
	_ "libdb.so/vm/programs/neofetch"
	_ "libdb.so/vm/programs/sixel"
	_ "libdb.so/vm/programs/spew"
	_ "libdb.so/vm/programs/termio"
)

const InitialCwd = "/"

//go:embed shellrc
var RC string

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
	var b bytes.Buffer
	b.Grow(1024)

	return func(env vm.Environment) string {
		q := env.Terminal.Query()

		b.Reset()
		b.WriteByte('\n')

		line1 := fmt.Sprintf("$ libdb.so @ %s", env.Cwd)
		lineprompt.Blend(&b, line1, q.Width, transBlend, lineprompt.Opts{
			LOD:       15,
			Underline: true,
		})

		b.WriteByte('\n')
		b.WriteString("\033[38;2;85;205;252m―❤―\033[0m\033[38;2;247;157;208m▶\033[m ")
		return b.String()
	}
}

func rgb(r, g, b uint8) colorful.Color {
	return colorful.Color{
		R: float64(r) / 255,
		G: float64(g) / 255,
		B: float64(b) / 255,
	}
}
