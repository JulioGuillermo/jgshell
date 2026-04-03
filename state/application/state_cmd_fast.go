package stateapplication

import (
	"strings"

	"github.com/acarl005/stripansi"
)

func (s *State) FastCmd(cmd string) (string, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	uuid := GetUUID()
	s.isRunning = false

	s.Clear()
	_, err := s.shell.Write([]byte(WrapCmd(cmd, uuid)))
	if err != nil {
		return err.Error(), -2
	}

	output := ""
	buf := make([]byte, 1024)
	started := false
	code := -10
	for {
		n, err := s.shell.Read(buf)
		if err != nil {
			break
		}
		if n <= 0 {
			continue
		}
		result := CleanOutput(output+string(buf[:n]), started, uuid)
		output = result.Output
		code = result.Code
		if result.Started {
			started = true
		}
		if !result.IsRunning {
			break
		}
	}

	return strings.TrimSpace(output), code
}

func (s *State) FastCmdClean(cmd string) (string, int) {
	output, code := s.FastCmd(cmd)
	cleanText := stripansi.Strip(output)
	return cleanText, code
}
