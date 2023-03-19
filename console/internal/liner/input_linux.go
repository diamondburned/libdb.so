package liner

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	getTermios = syscall.TCGETS
	setTermios = syscall.TCSETS
)

const (
	icrnl  = syscall.ICRNL
	inpck  = syscall.INPCK
	istrip = syscall.ISTRIP
	ixon   = syscall.IXON
	opost  = syscall.OPOST
	cs8    = syscall.CS8
	isig   = syscall.ISIG
	icanon = syscall.ICANON
	iexten = syscall.IEXTEN
)

type termios struct {
	syscall.Termios
}

const cursorColumn = false

// TerminalSupported returns true if the current terminal supports
// line editing features, and false if liner will use the 'dumb'
// fallback for input.
// Note that TerminalSupported does not check all factors that may
// cause liner to not fully support the terminal (such as stdin redirection)
func TerminalSupported() bool {
	bad := map[string]bool{"": true, "dumb": true, "cons25": true}
	return !bad[strings.ToLower(os.Getenv("TERM"))]
}

func initLinerTerminal(s *State) {
	if m, err := TerminalMode(); err == nil {
		s.origMode = *m.(*termios)
	} else {
		s.inputRedirected = true
	}
	if _, err := getMode(syscall.Stdout); err != nil {
		s.outputRedirected = true
	}
	if s.inputRedirected && s.outputRedirected {
		s.terminalSupported = false
	}
	if s.terminalSupported && !s.inputRedirected && !s.outputRedirected {
		mode := s.origMode
		mode.Iflag &^= icrnl | inpck | istrip | ixon
		mode.Cflag |= cs8
		mode.Lflag &^= syscall.ECHO | icanon | iexten
		mode.Cc[syscall.VMIN] = 1
		mode.Cc[syscall.VTIME] = 0
		mode.ApplyMode()

		winch := make(chan os.Signal, 1)
		signal.Notify(winch, syscall.SIGWINCH)
		s.winch = winch

		s.checkOutput()
	}

	if !s.outputRedirected {
		s.outputRedirected = !s.getColumns()
	}
}

func (s *State) supportedStartPrompt() {
	if s.terminalSupported {
		if m, err := TerminalMode(); err == nil {
			s.defaultMode = *m.(*termios)
			mode := s.defaultMode
			mode.Lflag &^= isig
			mode.ApplyMode()
		}
	}
}

func (s *State) supportedStopPrompt() {
	if s.terminalSupported {
		s.defaultMode.ApplyMode()
	}
}

func (s *State) close() {
	if !s.inputRedirected {
		s.origMode.ApplyMode()
	}
}
