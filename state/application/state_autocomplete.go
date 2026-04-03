package stateapplication

import (
	"fmt"
	"strings"

	"github.com/julioguillermo/jgshell/scripts"
)

type autocompleteConfig struct {
	interpreter string
	script      string
}

var autocompleteScripts = map[string]autocompleteConfig{
	"bash":       {"bash", "autocomplete/autocomplete_bash.sh"},
	"zsh":        {"zsh", "autocomplete/autocomplete_zsh.sh"},
	"fish":       {"fish", "autocomplete/autocomplete_fish.fish"},
	"powershell": {"pwsh", "autocomplete/autocomplete_powershell.ps1"},
	"nushell":    {"nu", "autocomplete/autocomplete_nushell.nu"},
	"sh":         {"sh", "autocomplete/autocomplete_sh.sh"},
}

// shellEscape escapes a string for safe use as a shell argument
func shellEscape(s string) string {
	escaped := strings.ReplaceAll(s, "'", "'\\''")
	return "'" + escaped + "'"
}

func (s *State) GetAutoComplete(line string, cursor int) []string {
	shellType := s.GetShell()

	cfg, ok := autocompleteScripts[shellType]
	if !ok {
		cfg = autocompleteScripts["bash"]
	}

	bytes, err := scripts.AutoCompleteScript.ReadFile(cfg.script)
	if err != nil {
		return []string{}
	}

	script := string(bytes)

	// Build command: pipe embedded script to interpreter with args
	var cmd string
	switch shellType {
	case "nushell":
		cmd = fmt.Sprintf("JG_LINE=%s JG_CURSOR=%d %s -c %s",
			shellEscape(line), cursor, cfg.interpreter, shellEscape(script))
	case "powershell":
		// PowerShell: write script to temp file and execute
		cmd = fmt.Sprintf("$f=[System.IO.Path]::GetTempFileName()+'.ps1';"+
			"[System.IO.File]::WriteAllText($f,@'\n%s\n'@);"+
			" %s -File $f -Line %s -Cursor %d;"+
			" Remove-Item $f -ErrorAction SilentlyContinue",
			script, cfg.interpreter, shellEscape(line), cursor)
	default:
		// bash/zsh/fish: write to temp file and execute
		cmd = fmt.Sprintf("f=$(mktemp); cat > \"$f\" <<'JGEOF'\n%s\nJGEOF\n"+
			"%s \"$f\" %s %d; rm -f \"$f\"",
			script, cfg.interpreter, shellEscape(line), cursor)
	}

	output, code := s.FastCmdClean(cmd)
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
