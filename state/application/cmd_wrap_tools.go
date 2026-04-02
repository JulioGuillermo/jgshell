package stateapplication

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var shellStartMarkerRegex = regexp.MustCompile(`\x1b]123;START\x07([^\s]+) ([^\s\n]+) >>>`)
var shellSimpleStartMarkerRegex = regexp.MustCompile(`\x1b]123;START\x07`)

var shellEndMarkerRegex = regexp.MustCompile(`\x1b]123;(\d+);DONE\x07`)

type CleanOutputResult struct {
	Output    string
	Started   bool
	IsRunning bool
	Username  string
	Pwd       string
	Code      int
	EndTime   *time.Time
}

func CleanOutput(output string) *CleanOutputResult {
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

func CleanSimpleOutput(output string) (string, bool, int) {
	loc := shellSimpleStartMarkerRegex.FindStringIndex(output)
	if loc != nil {
		output = output[loc[1]:]
	}

	match := shellEndMarkerRegex.FindStringSubmatch(output)
	if len(match) <= 1 {
		return output, false, -10
	}

	code, _ := strconv.Atoi(match[1])

	loc = shellEndMarkerRegex.FindStringIndex(output)
	if loc == nil {
		return output, false, code
	}

	output = output[:loc[0]]
	return output, true, code
}

func WrapCmd(message string) string {
	return fmt.Sprintf("printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;DONE\\007\\n\"\n", message)
}

func WrapSimpleCmd(message string) string {
	return fmt.Sprintf("printf \"\\033]123;START\\007\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;DONE\\007\\n\"\n", message)
}
