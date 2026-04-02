package stateapplication

func (s *State) FastCmd(cmd string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isRunning = false

	output := ""
	buf := make([]byte, 1024)
	end := false
	for !end {
		n, err := s.shell.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			output, end = CleanSimpleOutput(output + string(buf[:n]))
		}
	}
	return output
}
