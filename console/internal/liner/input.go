//go:build linux || js
// +build linux js

package liner

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type nexter struct {
	r   rune
	err error
}

// State represents an open terminal
type State struct {
	commonState
	origMode    termios
	defaultMode termios
	next        <-chan nexter
	winch       chan os.Signal
	pending     []rune
	useCHA      bool
	getWinSize  func() (uint16, uint16, bool)
}

// NewState initializes a new *State, and sets the terminal into raw mode. To
// restore the terminal to its previous state, call State.Close().
func NewState() *State {
	return NewStateStdin(os.Stdin, nil)
}

// NewStateStdin initializes a new *State with the given stdin.
func NewStateStdin(stdin io.Reader, winszFunc func() (uint16, uint16, bool)) *State {
	var s State
	s.r = bufio.NewReader(stdin)
	s.getWinSize = winszFunc

	s.terminalSupported = TerminalSupported()
	initLinerTerminal(&s)

	return &s
}

var errTimedOut = errors.New("timeout")

func (s *State) startPrompt() {
	s.supportedStartPrompt()
	s.restartPrompt()
}

func (s *State) inputWaiting() bool {
	return len(s.next) > 0
}

func (s *State) restartPrompt() {
	next := make(chan nexter, 200)
	go func() {
		for {
			var n nexter
			n.r, _, n.err = s.r.ReadRune()
			next <- n
			// Shut down nexter loop when an end condition has been reached
			if n.err != nil || n.r == '\n' || n.r == '\r' || n.r == ctrlC || n.r == ctrlD {
				close(next)
				return
			}
		}
	}()
	s.next = next
}

func (s *State) stopPrompt() {
	s.supportedStopPrompt()
}

func (s *State) nextPending(timeout <-chan time.Time) (rune, error) {
	select {
	case thing, ok := <-s.next:
		if !ok {
			return 0, ErrInternal
		}
		if thing.err != nil {
			return 0, thing.err
		}
		s.pending = append(s.pending, thing.r)
		return thing.r, nil
	case <-timeout:
		rv := s.pending[0]
		s.pending = s.pending[1:]
		return rv, errTimedOut
	}
}

func (s *State) readNext() (interface{}, error) {
	if len(s.pending) > 0 {
		rv := s.pending[0]
		s.pending = s.pending[1:]
		return rv, nil
	}
	var r rune
	select {
	case thing, ok := <-s.next:
		if !ok {
			return 0, ErrInternal
		}
		if thing.err != nil {
			return nil, thing.err
		}
		r = thing.r
	case <-s.winch:
		s.getColumns()
		return winch, nil
	}
	if r != esc {
		return r, nil
	}
	s.pending = append(s.pending, r)

	// Wait at most 50 ms for the rest of the escape sequence
	// If nothing else arrives, it was an actual press of the esc key
	timeout := time.After(50 * time.Millisecond)
	flag, err := s.nextPending(timeout)
	if err != nil {
		if err == errTimedOut {
			return flag, nil
		}
		return unknown, err
	}

	switch flag {
	case '[':
		code, err := s.nextPending(timeout)
		if err != nil {
			if err == errTimedOut {
				return code, nil
			}
			return unknown, err
		}
		switch code {
		case 'A':
			s.pending = s.pending[:0] // escape code complete
			return up, nil
		case 'B':
			s.pending = s.pending[:0] // escape code complete
			return down, nil
		case 'C':
			s.pending = s.pending[:0] // escape code complete
			return right, nil
		case 'D':
			s.pending = s.pending[:0] // escape code complete
			return left, nil
		case 'F':
			s.pending = s.pending[:0] // escape code complete
			return end, nil
		case 'H':
			s.pending = s.pending[:0] // escape code complete
			return home, nil
		case 'Z':
			s.pending = s.pending[:0] // escape code complete
			return shiftTab, nil
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			num := []rune{code}
			for {
				code, err := s.nextPending(timeout)
				if err != nil {
					if err == errTimedOut {
						return code, nil
					}
					return nil, err
				}
				switch code {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					num = append(num, code)
				case ';':
					// Modifier code to follow
					// This only supports Ctrl-left and Ctrl-right for now
					x, _ := strconv.ParseInt(string(num), 10, 32)
					if x != 1 {
						// Can't be left or right
						rv := s.pending[0]
						s.pending = s.pending[1:]
						return rv, nil
					}
					num = num[:0]
					for {
						code, err = s.nextPending(timeout)
						if err != nil {
							if err == errTimedOut {
								rv := s.pending[0]
								s.pending = s.pending[1:]
								return rv, nil
							}
							return nil, err
						}
						switch code {
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							num = append(num, code)
						case 'C', 'D':
							// right, left
							mod, _ := strconv.ParseInt(string(num), 10, 32)
							if mod != 5 {
								// Not bare Ctrl
								rv := s.pending[0]
								s.pending = s.pending[1:]
								return rv, nil
							}
							s.pending = s.pending[:0] // escape code complete
							if code == 'C' {
								return wordRight, nil
							}
							return wordLeft, nil
						default:
							// Not left or right
							rv := s.pending[0]
							s.pending = s.pending[1:]
							return rv, nil
						}
					}
				case '~':
					s.pending = s.pending[:0] // escape code complete
					x, _ := strconv.ParseInt(string(num), 10, 32)
					switch x {
					case 2:
						return insert, nil
					case 3:
						return del, nil
					case 5:
						return pageUp, nil
					case 6:
						return pageDown, nil
					case 1, 7:
						return home, nil
					case 4, 8:
						return end, nil
					case 15:
						return f5, nil
					case 17:
						return f6, nil
					case 18:
						return f7, nil
					case 19:
						return f8, nil
					case 20:
						return f9, nil
					case 21:
						return f10, nil
					case 23:
						return f11, nil
					case 24:
						return f12, nil
					default:
						return unknown, nil
					}
				default:
					// unrecognized escape code
					rv := s.pending[0]
					s.pending = s.pending[1:]
					return rv, nil
				}
			}
		}

	case 'O':
		code, err := s.nextPending(timeout)
		if err != nil {
			if err == errTimedOut {
				return code, nil
			}
			return nil, err
		}
		s.pending = s.pending[:0] // escape code complete
		switch code {
		case 'c':
			return wordRight, nil
		case 'd':
			return wordLeft, nil
		case 'H':
			return home, nil
		case 'F':
			return end, nil
		case 'P':
			return f1, nil
		case 'Q':
			return f2, nil
		case 'R':
			return f3, nil
		case 'S':
			return f4, nil
		default:
			return unknown, nil
		}
	case 'b':
		s.pending = s.pending[:0] // escape code complete
		return altB, nil
	case 'd':
		s.pending = s.pending[:0] // escape code complete
		return altD, nil
	case bs:
		s.pending = s.pending[:0] // escape code complete
		return altBs, nil
	case 'f':
		s.pending = s.pending[:0] // escape code complete
		return altF, nil
	case 'y':
		s.pending = s.pending[:0] // escape code complete
		return altY, nil
	default:
		rv := s.pending[0]
		s.pending = s.pending[1:]
		return rv, nil
	}

	// not reached
	return r, nil
}

// Close returns the terminal to its previous mode
func (s *State) Close() error {
	signal.Stop(s.winch)
	s.close()
	return nil
}
