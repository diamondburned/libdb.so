package ansi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/mattn/go-runewidth"
	"libdb.so/vm"
)

const (
	EchoOff = "\x1b[?25l"
	EchoOn  = "\x1b[?25h"

	ClearLine    = "\x1b[2K"
	ClearToStart = "\x1b[1K"
	ClearToEnd   = "\x1b[K"

	MoveCursorToStart = "\x1b[1G"
	MoveCursorToEnd   = "\x1b[G"
)

func MoveCursorUp(n int) string {
	return fmt.Sprintf("\x1b[%dA", n)
}

func MoveCursorDown(n int) string {
	return fmt.Sprintf("\x1b[%dB", n)
}

func MoveCursorRight(n int) string {
	return fmt.Sprintf("\x1b[%dC", n)
}

func MoveCursorLeft(n int) string {
	return fmt.Sprintf("\x1b[%dD", n)
}

var ansiLinkRe = regexp.MustCompile(`(?m)\x1b]8;;([^\x1b]*)\x1b\\([^\x1b]*)\x1b]8;;\x1b\\`)

const ansiLinkf = "\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\"

func init() {
	// keep in sync with the regex
	if !ansiLinkRe.MatchString(ansiLinkf) {
		panic("ansiLinkf is not in sync with ansiLinkRe")
	}
}

// Link formats a link with the given text and url.
func Link(text, url string) string {
	return fmt.Sprintf(ansiLinkf, url, text)
}

// StringWidth returns the width of the given string, ignoring ANSI escape
// sequences.
func StringWidth(str string) int {
	// Undo all the ANSI links.
	str = ansiLinkRe.ReplaceAllString(str, "$2")
	return ansiLength(str)
}

func ansiLength(str string) int {
	if str == "" {
		return 0
	}

	parsed, err := ansi.Parse(str)
	if err != nil {
		return 0
	}

	var result int
	for _, element := range parsed {
		result += runewidth.StringWidth(element.Label)
	}

	return result
}

// PrintAligned prints three strings, left-aligned, center-aligned, and
// right-aligned, with the given width. If width is 0, the terminal width is
// used.
func PrintAligned(env vm.Environment, width int, left, center, right string) {
	w := env.Terminal.Stdout
	if width == 0 {
		width = env.Terminal.Query().Width
	}

	leftWidth := StringWidth(left)
	rightWidth := StringWidth(right)
	centerWidth := StringWidth(center)
	halfCenterWidth := centerWidth / 2

	var totalWidth int
	if centerWidth > 0 {
		totalWidth = max(
			2*(leftWidth+halfCenterWidth+1),
			2*(halfCenterWidth+rightWidth+1),
		)
	} else {
		totalWidth = leftWidth + rightWidth
	}
	if totalWidth > width {
		centerPad := max((width-centerWidth)/2, 0)
		rightPad := max(width-rightWidth, 0)

		if left != "" {
			fmt.Fprintln(w, left)
		}

		if center != "" {
			fmt.Fprint(w, strings.Repeat(" ", centerPad))
			fmt.Fprintln(w, center)
		}

		if right != "" {
			fmt.Fprint(w, strings.Repeat(" ", rightPad))
			fmt.Fprintln(w, right)
		}

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
