package stateapplication

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed autocomplete.sh
var AutoCompleteScript embed.FS

func (s *State) GetAutoComplete(line string, cursor int) []string {
	bytes, err := StatusScript.ReadFile("status.sh")
	if err != nil {
		return []string{}
	}

	script := string(bytes)
	script = strings.ReplaceAll(script, "%{GO_LINE}", line)
	script = strings.ReplaceAll(script, "%{GO_CURSOR}", fmt.Sprint(cursor))

	output, code := s.FastCmdClean(script)
	if code != 0 {
		return []string{}
	}
	var completions []string
	for l := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			completions = append(completions, l)
		}
	}
	return completions
}
