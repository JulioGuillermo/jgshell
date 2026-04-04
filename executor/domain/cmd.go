package executordomain

import (
	"regexp"
	"time"

	"github.com/julioguillermo/jgshell/tools"
)

const (
	Start = `\x1b\[\?(1049|47|1047)h`
	End   = `\x1b\[\?(1049|47|1047)l`
)

type Cmd struct {
	UUID     string
	SH       string
	PWD      string
	USER     string
	Cmd      string
	Output   string
	Started  bool
	ExitCode int
	Error    error
	Start    *time.Time
	End      *time.Time

	cleanedOutput *string
}

func (c *Cmd) IsRunning() bool {
	return c.ExitCode == -10 && c.End == nil
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

func (c *Cmd) GetRenderableOutput() string {
	output := regexp.MustCompile("(?s)"+Start+".*?"+End).
		ReplaceAllString(c.Output, "")
	output = regexp.MustCompile("(?s)"+Start+".*").
		ReplaceAllString(output, "")
	return output
}

func (c *Cmd) CleanOuput() string {
	if !c.Started {
		return ""
	}
	if c.IsRunning() {
		return tools.CleanText(c.GetRenderableOutput())
	}
	if c.cleanedOutput == nil {
		c.cleanedOutput = new(string)
		*c.cleanedOutput = tools.CleanText(c.GetRenderableOutput())
	}
	return *c.cleanedOutput
}

func (c *Cmd) IsFullScreen() bool {
	if !c.Started {
		return false
	}
	if !c.IsRunning() {
		return false
	}

	idxs := regexp.MustCompile(Start).FindAllStringIndex(c.Output, -1)
	if len(idxs) == 0 {
		return false
	}

	fsInit := idxs[len(idxs)-1][0]
	if fsInit == -1 {
		return false
	}

	idxs = regexp.MustCompile(End).FindAllStringIndex(c.Output, -1)
	if len(idxs) == 0 {
		return true
	}
	fsEnd := idxs[len(idxs)-1][0]
	if fsInit == -1 {
		return true
	}

	return fsEnd < fsInit
}
