package stateapplication

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var shellStartMarkerRegex = regexp.MustCompile(`\033]123;START;([^\s]+);([^\s\n]+);>>>\007`)
var shellEndMarkerRegex = regexp.MustCompile(`\033]123;(\d+);DONE\007`)

type CleanOutputResult struct {
	Output    string
	Started   bool
	IsRunning bool
	Username  string
	Pwd       string
	Code      int
	EndTime   *time.Time
}

func CleanOutput(output string, started bool, uuid string) *CleanOutputResult {
	result := &CleanOutputResult{
		Output:    output,
		Started:   started,
		IsRunning: true,
		Code:      -10,
	}

	if !result.Started {
		match := shellStartMarkerRegex.FindStringSubmatch(result.Output)
		if len(match) > 2 {
			result.Started = true
			loc := shellStartMarkerRegex.FindStringIndex(result.Output)
			result.Output = result.Output[loc[1]:]
			result.Username = match[1]
			result.Pwd = match[2]
		}
		return result
	}

	// result.Output = shellStartMarkerRegex.ReplaceAllString(result.Output, "")

	match := shellEndMarkerRegex.FindStringSubmatch(result.Output)
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

func WrapCmd(message string, uuid string) string {
	message = strings.TrimSpace(message)
	// return fmt.Sprintf("printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;DONE\\007\\n\"\n", message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { \n%s\n } ; printf \"\\033]123;$?;$JG_BLOCK_UUID;DONE\\007\\n\"\n", uuid, message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; printf \"\\033]123;START\\007$(whoami) $(pwd) >>>\\n\" ; { %s }\n", uuid, message)
	// return fmt.Sprintf("export JG_BLOCK_UUID=%s; { \n%s\n }\n", uuid, message)
	return fmt.Sprintf("printf \"\\033]123;START;$(whoami);$(pwd);>>>\\007\" ; { %s ; }\n", message)
	// return "{\n" + strings.TrimSpace(message) + "\n}\n"
}

func GetUUID() string {
	return uuid.New().String()
}
