package wrapperapplication

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type CmdWrapper struct {
	StartMarkerRegex *regexp.Regexp
	EndMarkerRegex   *regexp.Regexp
}

func NewCmdWrapperWithMarkers(startMarker, endMarker *regexp.Regexp) *CmdWrapper {
	return &CmdWrapper{
		StartMarkerRegex: startMarker,
		EndMarkerRegex:   endMarker,
	}
}

func NewCmdWrapper() *CmdWrapper {
	return NewCmdWrapperWithMarkers(
		regexp.MustCompile(`\033]123;START;(.+);(.+);>>>\007`),
		regexp.MustCompile(`\033]123;(\d+);DONE\007`),
	)
}

func (w *CmdWrapper) WrapCmd(command string) string {
	return fmt.Sprintf("printf \"\\033]123;START;%%s;%%s;>>>\\007\" \"$(whoami)\" \"$(pwd)\" ; { %s ; }\n", command)
}

func (w *CmdWrapper) UnwrapCmd(output string, started bool) *wrapperdomain.CmdUnwrapResult {
	result := &wrapperdomain.CmdUnwrapResult{
		Output:    output,
		Started:   started,
		IsRunning: true,
		Code:      -10,
	}

	if !result.Started {
		w.processStart(result)
		return result
	}

	w.processEnd(result)
	return result
}

func (w *CmdWrapper) processStart(result *wrapperdomain.CmdUnwrapResult) {
	match := w.StartMarkerRegex.FindStringSubmatch(result.Output)
	if len(match) <= 2 {
		return
	}
	result.Started = true

	result.User = match[1]
	result.Pwd = match[2]

	loc := w.StartMarkerRegex.FindStringIndex(result.Output)
	result.Output = result.Output[loc[1]:]
}

func (w *CmdWrapper) processEnd(result *wrapperdomain.CmdUnwrapResult) {
	match := w.EndMarkerRegex.FindStringSubmatch(result.Output)
	if len(match) <= 1 {
		return
	}

	end := time.Now()
	result.EndTime = &end
	result.IsRunning = false
	result.Code, _ = strconv.Atoi(match[1])

	loc := w.EndMarkerRegex.FindStringIndex(result.Output)
	result.Output = result.Output[:loc[0]]
}
