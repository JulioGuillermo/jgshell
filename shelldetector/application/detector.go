package application

import (
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
	psVersion, err := d.executor.Run("printf \"$PSVersionTable\"")
	if err != nil {
		return "", err
	}
	if strings.Contains(psVersion, "PSVersionHashTable") {
		return "powershell", nil
	}
	return d.executor.Run("ps -p $$ -o comm= 2>/dev/null | sed 's/-//' || echo \"sh\"")
}
