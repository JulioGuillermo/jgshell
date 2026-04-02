package stateapplication

func (s *State) ShowInput() bool {
	return !s.isRunning
}

func (s *State) GetAutoComplete() []string {
	return []string{}
}
