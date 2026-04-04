package statusinfrastructure

import (
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	"github.com/julioguillermo/jgshell/scripts"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
	statusdomain "github.com/julioguillermo/jgshell/status/domain"
)

type StatusLoader struct {
	shellDetector     shelldetectordomain.ShellDetector
	shellFastExecutor executordomain.FastExecutor
}

func NewStatusLoader(shellDetector shelldetectordomain.ShellDetector, shellFastExecutor executordomain.FastExecutor) *StatusLoader {
	return &StatusLoader{
		shellDetector:     shellDetector,
		shellFastExecutor: shellFastExecutor,
	}
}

func (s *StatusLoader) Load() (*statusdomain.Status, error) {
	sh, err := s.shellDetector.DetectShell()
	if err != nil {
		return nil, err
	}
	result := &statusdomain.Status{
		Shell: sh,
	}

	script, err := s.getScript(sh)
	if err != nil {
		return result, err
	}

	output, _, err := s.shellFastExecutor.RunAndClean(script)
	if err != nil {
		return result, err
	}

	s.parseOS(result, output)
	s.parseUser(result, output)
	s.parseDir(result, output)
	s.parseGit(result, output)

	return result, nil
}

func (s *StatusLoader) getScript(sh string) (string, error) {
	if sh == "powershell" {
		script, err := scripts.StatusScript.ReadFile("status/status.ps1")
		if err != nil {
			return "", err
		}
		return string(script), nil
	}

	script, err := scripts.StatusScript.ReadFile("status/status.sh")
	if err != nil {
		return "", err
	}
	return string(script), nil
}
