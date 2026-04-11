package wrapperapplication

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type CmdWrapper struct {
	StartMarkerFastRegex *regexp.Regexp
	StartMarkerRegex     *regexp.Regexp
	EndMarkerRegex       *regexp.Regexp
}

func NewCmdWrapper() *CmdWrapper {
	return &CmdWrapper{
		StartMarkerFastRegex: regexp.MustCompile(wrapperdomain.REWrapStartFast),
		StartMarkerRegex:     regexp.MustCompile(wrapperdomain.REWrapStart),
		EndMarkerRegex:       regexp.MustCompile(wrapperdomain.REWrapDone),
	}
}

func (w *CmdWrapper) FastWrapCmd(sh, command string) string {
	return w.wrap(sh, command, true)
}

func (w *CmdWrapper) WrapCmd(sh, command string) string {
	return w.wrap(sh, command, false)
}

func (w *CmdWrapper) wrap(sh, command string, fast bool) string {
	var start string
	if fast {
		start = wrapperdomain.WrapperStartFast
	} else {
		start = wrapperdomain.WrapperStart
	}

	command = strings.TrimSpace(command)
	switch sh {
	case "powershell":
		return fmt.Sprintf("%s ; . {\n%s\n}\n\n\n", start, command)
	default:
		return fmt.Sprintf("%s ; {\n%s\n}\n", start, command)
	}
}

func (w *CmdWrapper) FastUnwrapCmd(output string, started bool) *wrapperdomain.CmdUnwrapResult {
	return w.unwrap(output, started, true)
}

func (w *CmdWrapper) UnwrapCmd(output string, started bool) *wrapperdomain.CmdUnwrapResult {
	return w.unwrap(output, started, false)
}

func (w *CmdWrapper) unwrap(output string, started bool, fast bool) *wrapperdomain.CmdUnwrapResult {
	result := &wrapperdomain.CmdUnwrapResult{
		Output:    output,
		Started:   started,
		IsRunning: true,
		Code:      -10,
	}

	if !result.Started {
		w.processStart(result, fast)
		if !result.Started {
			return result
		}
	}

	w.processEnd(result)
	return result
}

func (w *CmdWrapper) processStart(result *wrapperdomain.CmdUnwrapResult, fast bool) {
	var startRegex *regexp.Regexp
	if fast {
		startRegex = w.StartMarkerFastRegex
	} else {
		startRegex = w.StartMarkerRegex
	}

	match := startRegex.FindStringSubmatch(result.Output)
	if len(match) <= 2 {
		return
	}
	result.Started = true

	result.User = match[1]
	result.Pwd = match[2]

	loc := startRegex.FindStringIndex(result.Output)
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
