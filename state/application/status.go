package stateapplication

import statedomain "github.com/julioguillermo/jgshell/state/domain"

type Status struct {
	user  string
	dir   string
	os    string
	shell string
	git   *statedomain.GitStats
}

func (s *Status) User() string {
	return s.user
}

func (s *Status) Dir() string {
	return s.dir
}

func (s *Status) OS() string {
	return s.os
}

func (s *Status) Shell() string {
	return s.shell
}

func (s *Status) Git() *statedomain.GitStats {
	return s.git
}
