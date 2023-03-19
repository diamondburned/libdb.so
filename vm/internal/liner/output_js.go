package liner

const cursorColumn = false // unsure

func (s *State) getColumns() bool {
	if s.getWinSize == nil {
		return false
	}

	_, col, ok := s.getWinSize()
	if !ok {
		return false
	}

	s.columns = int(col)
	return true
}
