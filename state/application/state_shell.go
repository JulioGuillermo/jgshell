package stateapplication

import (
	"fmt"
	"strings"
)

func (s *State) GetSimpleOutput(cmd string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	uuid := GetUUID()
	s.isRunning = false

	fmt.Fprintf(
		s.shell,
		`printf "<<<JGSHELL_START;%%s>>> %%s <<<JGSHELL_END;%%s>>>\n" "%s" "$(%s)" "%s"
`,
		uuid, cmd, uuid,
	)

	Start := `<<<JGSHELL_START;` + uuid + `>>>`
	End := `<<<JGSHELL_END;` + uuid + `>>>`

	buf := make([]byte, 1024)
	output := ""
	for {
		n, err := s.shell.Read(buf)
		if err != nil {
			break
		}
		if n <= 0 {
			continue
		}

		output += string(buf[:n])
		// output = stripansi.Strip(output)
		if strings.Contains(output, End) {
			break
		}
	}

	idx := strings.Index(output, Start)
	if idx != -1 {
		output = output[idx+len(Start):]
	}
	idx = strings.Index(output, End)
	if idx != -1 {
		output = output[:idx]
	}

	return strings.TrimSpace(output)
}

func (s *State) GetShell() string {
	psVersion := s.GetSimpleOutput("printf \"$PSVersionTable\"")
	if strings.Contains(psVersion, "PSVersionHashTable") {
		return "powershell"
	}
	return s.GetSimpleOutput("ps -p $$ -o comm= 2>/dev/null | sed 's/-//' || echo \"sh\"")
}
