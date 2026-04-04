package autocompleteinfrastructure

import (
	"fmt"
	"strings"
)

func (a *Autocomplete) buildCmd(cursor int, line, shell, cmd, script string) string {
	switch shell {
	case "nushell":
		return fmt.Sprintf("JG_LINE=%s JG_CURSOR=%d %s -c %s",
			a.shellEscape(line), cursor, cmd, a.shellEscape(script))
	case "powershell":
		return fmt.Sprintf("$f=[System.IO.Path]::GetTempFileName()+'.ps1';"+
			"[System.IO.File]::WriteAllText($f,@'\n%s\n'@);"+
			" %s -File $f -Line %s -Cursor %d;"+
			" Remove-Item $f -ErrorAction SilentlyContinue",
			script, cmd, a.shellEscape(line), cursor)
	default:
		return fmt.Sprintf("f=$(mktemp); cat > \"$f\" <<'JGEOF'\n%s\nJGEOF\n"+
			"%s \"$f\" %s %d; rm -f \"$f\"",
			script, cmd, a.shellEscape(line), cursor)
	}
}

func (a *Autocomplete) shellEscape(s string) string {
	// TODO: Check this
	escaped := strings.ReplaceAll(s, "'", "'\\''")
	return "'" + escaped + "'"
}
