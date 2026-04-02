package stateapplication

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var shellEndMarkerRegex = regexp.MustCompile(`\x1b]123;(\d+);DONE\x07`)
var shellStartMarkerRegex = regexp.MustCompile(`\x1b]123;START\x07([^\s]+) ([^\s\n]+) >>>`)

type CleanOutputResult struct {
	Output    string
	Started   bool
	IsRunning bool
	Username  string
	Pwd       string
	Code      int
	EndTime   *time.Time
}

func cleanOutput(output string) *CleanOutputResult {
	result := &CleanOutputResult{
		Output:    output,
		IsRunning: true,
		Code:      -10,
	}

	match := shellStartMarkerRegex.FindStringSubmatch(result.Output)
	if len(match) > 2 {
		result.Started = true
		loc := shellStartMarkerRegex.FindStringIndex(result.Output)
		result.Output = result.Output[loc[1]:]
		result.Username = match[1]
		result.Pwd = match[2]
	}

	match = shellEndMarkerRegex.FindStringSubmatch(result.Output)
	if len(match) <= 1 {
		return result
	}

	result.Code, _ = strconv.Atoi(match[1])
	loc := shellEndMarkerRegex.FindStringIndex(result.Output)

	result.Output = result.Output[:loc[0]]

	end := time.Now()
	result.EndTime = &end
	result.IsRunning = false
	return result
}

func wrapCmd(message string) string {
	return fmt.Sprintf("printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { %s ; } ; printf \"\\033]123;$?;DONE\\007\\n\"\n", message)
}
