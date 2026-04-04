package application

import (
	"fmt"

	"github.com/julioguillermo/jgshell/scripts"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
)

type ShellWrapper struct {
	shell         shelldomain.Shell
	shellDetector shelldetectordomain.ShellDetector
}

func NewShellWrapper(shell shelldomain.Shell, shellDetector shelldetectordomain.ShellDetector) *ShellWrapper {
	return &ShellWrapper{
		shell:         shell,
		shellDetector: shellDetector,
	}
}

func (s *ShellWrapper) WrapShell() error {
	sh, err := s.shellDetector.DetectShell()
	if err != nil {
		return err
	}

	script, source := s.getWrapper(sh)
	if script == "" {
		return fmt.Errorf("Not wrapper script for shell: %s", sh)
	}

	if !source {
		s.shell.Write([]byte(script))
		return nil
	}

	if sh == "powershell" {
		return s.pwshWrap(script)
	}

	return s.shWrap(script)
}

func (s *ShellWrapper) getWrapper(sh string) (string, bool) {
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

func (s *ShellWrapper) pwshWrap(script string) error {
	_, err := fmt.Fprintf(
		s.shell,
		`$setupScript = @'
%s
'@
# Ejecuta el contenido del string en la sesión actual
Invoke-Expression $setupScript

Write-Output "Wrapper script executed successfully"
`,
		script,
	)
	return err
}

func (s *ShellWrapper) shFileWrap(script string) error {
	_, err := fmt.Fprintf(
		s.shell,
		`cat << 'EOF' > .wrapper
%s
EOF

chmod +x .wrapper
. ./.wrapper

printf "Wrapper script executed successfully"
`,
		script,
	)
	return err
}

func (s *ShellWrapper) shWrap(script string) error {
	_, err := fmt.Fprintf(
		s.shell,
		`eval "$(cat << 'EOF'
%s
EOF
)"
printf "Wrapper script executed successfully\n"
`,
		script,
	)
	return err
}
