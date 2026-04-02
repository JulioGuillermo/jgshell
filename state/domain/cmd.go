package statedomain

import (
	"time"

	"github.com/julioguillermo/jgshell/tools"
)

type Cmd struct {
	Cmd      string
	Output   string
	PWD      string
	USER     string
	Started  bool
	ExitCode int
	Error    error
	Start    *time.Time
	End      *time.Time
}

func NewCmd(cmd string) *Cmd {
	return &Cmd{
		Cmd: cmd,
	}
}

func (c *Cmd) GetDuration() time.Duration {
	if c.Start == nil {
		return 0
	}
	if c.End == nil {
		return time.Since(*c.Start)
	}
	return c.End.Sub(*c.Start)
}

func (c *Cmd) CleanOuput() string {
	if !c.Started {
		return ""
	}
	return tools.CleanText(c.Output)
}
