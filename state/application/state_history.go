package stateapplication

import statedomain "github.com/julioguillermo/jgshell/state/domain"

func (s *State) GetHistory() []statedomain.Cmd {
	return s.history
}

func (s *State) lastCmd() *statedomain.Cmd {
	if len(s.history) == 0 {
		return nil
	}
	return &s.history[len(s.history)-1]
}
