package stateapplication

type Status struct {
	user  string
	dir   string
	os    string
	shell string
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
