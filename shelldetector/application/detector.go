package shelldetectorapplication

import (
	"strconv"
	"strings"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
)

type ShellDetector struct {
	executor executordomain.SimpleExecutor
}

func NewShellDetector(executor executordomain.SimpleExecutor) *ShellDetector {
	return &ShellDetector{
		executor: executor,
	}
}

func (d *ShellDetector) DetectShell() (string, error) {
	pwsh, err := d.isPWSH()
	if err != nil {
		return "", err
	}
	if pwsh {
		return "powershell", nil
	}
	fish, err := d.isFish()
	if err != nil {
		return "", err
	}
	if fish {
		return "fish", nil
	}
	return d.executor.Run("ps -p $$ -o comm= 2>/dev/null | sed 's/-//' || echo \"sh\"")
}

func (d *ShellDetector) isPWSH() (bool, error) {
	psVersion, err := d.executor.Run("printf \"$PSVersionTable\"")
	if err != nil {
		return false, err
	}
	return strings.Contains(psVersion, "PSVersionHashTable"), nil
}

func (d *ShellDetector) isFish() (bool, error) {
	fish, err := d.executor.Run("printf \"$fish_pid\"")
	if err != nil {
		return false, err
	}
	if _, err := strconv.Atoi(fish); err == nil {
		return true, nil
	}
	return false, nil
}
