package stateapplication

func (s *State) ShowStatusBar() bool {
	return true
}

func (s *State) GetDir() string {
	return "/home"
}

func (s *State) GetTime() string {
	return ""
}
