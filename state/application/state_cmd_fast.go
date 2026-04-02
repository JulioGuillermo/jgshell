package stateapplication

import (
	"strings"

	"github.com/acarl005/stripansi"
)

func (s *State) FastCmd(cmd string) (string, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isRunning = false

	cmd = WrapSimpleCmd(cmd)
	s.shell.Write([]byte(cmd))

	output := ""
	buf := make([]byte, 1024)
	end := false
	code := -10
	for !end {
		n, err := s.shell.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			output, end, code = CleanSimpleOutput(output + string(buf[:n]))
		}
	}

	return strings.TrimSpace(output), code
}

func (s *State) FastCmdClean(cmd string) (string, int) {
	output, code := s.FastCmd(cmd)
	cleanText := stripansi.Strip(output)
	return cleanText, code
}
