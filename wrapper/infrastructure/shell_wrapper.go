package wrapperinfrastructure

import (
	"fmt"

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

	config, err := s.getConfig(sh)
	if err != nil {
		return err
	}
	if config == nil {
		return fmt.Errorf("Not wrapper script for shell: %s", sh)
	}

	_, err = fmt.Fprintf(
		s.shell,
		config.Loader,
		config.Script,
	)
	return err
}
