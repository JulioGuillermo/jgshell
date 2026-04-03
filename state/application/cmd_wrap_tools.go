package stateapplication

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var shellStartMarkerRegex = regexp.MustCompile(`\x1b]123;START\x07([^\s]+) ([^\s\n]+) >>>`)
var shellSimpleStartMarkerRegex = regexp.MustCompile(`\x1b]123;START\x07`)

const (
	shellEndMarkerFormat = `\x1b]123;(\d+);%{UUID};DONE\x07`
)

type CleanOutputResult struct {
	Output    string
	Started   bool
	IsRunning bool
	Username  string
	Pwd       string
	Code      int
	EndTime   *time.Time
}

func CleanOutput(output string, uuid string) *CleanOutputResult {
	var shellEndMarkerRegex = regexp.MustCompile(strings.ReplaceAll(shellEndMarkerFormat, "%{UUID}", uuid))

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

func CleanSimpleOutput(output string, uuid string) (string, bool, int) {
	var shellEndMarkerRegex = regexp.MustCompile(strings.ReplaceAll(shellEndMarkerFormat, "%{UUID}", uuid))

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

func GetPS1() string {
	return "PS1=$(printf \"\\033]123;$?;$JG_BLOCK_UUID;DONE\\007\\n\")\n"
}

func WrapCmd(message string, uuid string) string {
	// return fmt.Sprintf("printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;DONE\\007\\n\"\n", message)
	return fmt.Sprintf("export JG_BLOCK_UUID=%s; printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;$JG_BLOCK_UUID;DONE\\007\\n\"\n", uuid, message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { %s }\n", uuid, message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; { \n%s\n }\n", uuid, message)
}

func WrapSimpleCmd(message string, uuid string) string {
	// return fmt.Sprintf("printf \"\\033]123;START\\007\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;DONE\\007\\n\"\n", message)
	return fmt.Sprintf("export JG_BLOCK_UUID=%s; printf \"\\033]123;START\\007\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;$JG_BLOCK_UUID;DONE\\007\\n\"\n", uuid, message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; printf \"\\033]123;START\\007\\n\" ; { \n%s\n }\n", uuid, message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; { \n%s\n }\n", uuid, message)
}

func GetUUID() string {
	return uuid.New().String()
}

// GetWrapperScriptName returns the wrapper filename for a given shell type
func GetWrapperScriptName(shellType string) string {
	switch shellType {
	case "bash":
		return "wrapper_bash"
	case "zsh":
		return "wrapper_zsh"
	case "fish":
		return "wrapper_fish"
	case "nushell":
		return "wrapper_nushell"
	case "powershell":
		return "wrapper_powershell"
	default:
		return "wrapper_sh"
	}
}

// GetWrapperSourceCmd returns the shell command to source the wrapper script
func GetWrapperSourceCmd(wrapperPath string, shellType string) string {
	switch shellType {
	case "fish":
		return fmt.Sprintf("source %s", wrapperPath)
	case "powershell":
		return fmt.Sprintf(". %s", wrapperPath)
	case "nushell":
		return fmt.Sprintf("source %s", wrapperPath)
	default:
		return fmt.Sprintf(". %s", wrapperPath)
	}
}
