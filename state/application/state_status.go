package stateapplication

import (
	"github.com/julioguillermo/jgshell/scripts"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func (s *State) lastCmd() *statedomain.Cmd {
	if len(s.history) == 0 {
		return nil
	}
	return &s.history[len(s.history)-1]
}

func (s *State) GetHistory() []statedomain.Cmd {
	return s.history
}

func (s *State) IsRunning() bool {
	return s.isRunning
}

func (s *State) GetShell() string {
	script, err := scripts.ShellScript.ReadFile("shell/shell.sh")
	if err != nil {
		return "UNKNOWN"
	}
	output, _ := s.FastCmdClean(string(script))
	return output
}

func (s *State) GetStatus() statedomain.Status {
	status := &Status{}
	status.Load(s)
	return status
}
