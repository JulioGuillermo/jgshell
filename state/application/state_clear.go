package stateapplication

func (s *State) Clear() {
	buf := make([]byte, 2048)
	for {
		n, err := s.shell.Read(buf)
		if err != nil {
			break
		}
		if n < 1024 {
			break
		}
	}
}
