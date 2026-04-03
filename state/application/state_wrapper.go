package stateapplication

import (
	"fmt"

	"github.com/julioguillermo/jgshell/scripts"
)

func getWrapper(sh string) (string, bool) {
	var err error
	var bytes []byte
	var source bool
	switch sh {
	case "powershell":
		bytes, err = scripts.WrapperScript.ReadFile("wrapper/wrapper_powershell")
	case "bash":
		source = true
		bytes, err = scripts.WrapperScript.ReadFile("wrapper/wrapper_bash")
	case "zsh":
		source = true
		bytes, err = scripts.WrapperScript.ReadFile("wrapper/wrapper_zsh")
	case "fish":
		source = true
		bytes, err = scripts.WrapperScript.ReadFile("wrapper/wrapper_fish")
	case "nu":
		bytes, err = scripts.WrapperScript.ReadFile("wrapper/wrapper_nu")
	default:
		source = true
		bytes, err = scripts.WrapperScript.ReadFile("wrapper/wrapper_sh")
	}

	if err != nil {
		return "", false
	}
	return string(bytes), source
}

func (s *State) Wrap() {
	sh := s.GetShell()
	script, source := getWrapper(sh)
	if script == "" {
		return
	}

	if !source {
		s.shell.Write([]byte(script))
		return
	}

	fmt.Fprintf(
		s.shell,
		`cat << 'EOF' > .wrapper.sh
%s
EOF

chmod +x .wrapper.sh
. ./.wrapper.sh

printf "Wrapper script executed successfully"
`,
		script)

	buf := make([]byte, 1024)
	for {
		n, err := s.shell.Read(buf)
		if err != nil {
			break
		}
		if n < 100 {
			break
		}
	}
}
