package ansi

import (
	"fmt"
	"regexp"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/mattn/go-runewidth"
)

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
