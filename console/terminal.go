package console

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
	"libdb.so/console/internal/syncg"
)

// IO defines our needed IO files.
type IO struct {
	Stdin  io.ReadCloser
	Stdout io.Writer
	Stderr io.Writer
}

func (io IO) makeRaw() (func() error, error) {
	stdin, ok := io.Stdin.(*os.File)
	if !ok {
		return nil, fmt.Errorf("stdin is not a file but %T", io.Stdin)
	}

	oldState, err := terminal.MakeRaw(int(stdin.Fd()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to put terminal in raw mode")
	}

	return func() error {
		return terminal.Restore(0, oldState)
	}, nil
}

// TerminalQuery is a query to the terminal. It contains relevant terminal info
// needed for various purposes. For the most part, it's an extension to
// TIOCGWINSZ.
type TerminalQuery struct {
	// Width is the width of the terminal.
	Width int
	// Height is the height of the terminal.
	Height int
	// XPixel is the width of the terminal in pixels.
	XPixel int
	// YPixel is the height of the terminal in pixels.
	YPixel int
	// SIXEL is true if the terminal supports SIXEL.
	SIXEL bool
}

type terminalQueryUpdater struct {
	current syncg.AtomicValue[TerminalQuery]
	subs    syncg.Map[chan<- TerminalQuery, struct{}]
}

func (u *terminalQueryUpdater) set(q TerminalQuery) {
	u.current.Store(q)
	u.subs.Range(func(ch chan<- TerminalQuery, _ struct{}) bool {
		select {
		case ch <- q:
		default:
		}
		return true
	})
}

func (u *terminalQueryUpdater) subscribe(ch chan<- TerminalQuery) {
	u.subs.Store(ch, struct{}{})
}

func (u *terminalQueryUpdater) unsubscribe(ch chan<- TerminalQuery) {
	u.subs.Delete(ch)
}

// Terminal describes the current terminal state. It is mostly used by the
// programs.
type Terminal struct {
	IO
	query *terminalQueryUpdater
}

// NewTerminal creates a new terminal.
func NewTerminal(io IO, query TerminalQuery) *Terminal {
	var q terminalQueryUpdater
	q.set(query)
	return &Terminal{
		IO:    io,
		query: &q,
	}
}

// Query returns the current terminal query.
func (t *Terminal) Query() TerminalQuery {
	q, _ := t.query.current.Load()
	return q
}

// UpdateQuery updates the terminal query.
func (t *Terminal) UpdateQuery(q TerminalQuery) {
	t.query.set(q)
}

// Subscribe subscribes the given channel to terminal queries. It returns a
// function that, when called, unsubscribes the given channel from terminal
// queries. Queries that cannot be sent to the channel are dropped.
func (t *Terminal) Subscribe(ch chan<- TerminalQuery) func() {
	t.query.subscribe(ch)
	return func() { t.query.unsubscribe(ch) }
}

// Write writes to the terminal's stdout.
func (t *Terminal) Write(b []byte) (int, error) {
	return t.Stdout.Write(b)
}

// Read reads from the terminal's stdin.
func (t *Terminal) Read(b []byte) (int, error) {
	return t.Stdin.Read(b)
}
