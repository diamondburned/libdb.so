package liner

type termios = noopMode

func TerminalSupported() bool {
	return true // no stdin on js anyway
}

type noopMode struct{}

func (n noopMode) ApplyMode() error {
	return nil
}

// TerminalMode returns a noop InputModeSetter on this platform.
func TerminalMode() (ModeApplier, error) {
	return noopMode{}, nil
}

func initLinerTerminal(s *State) {
	s.inputRedirected = false
	s.outputRedirected = false
}

func (s *State) supportedStartPrompt() {}

func (s *State) supportedStopPrompt() {}

func (s *State) close() {}
