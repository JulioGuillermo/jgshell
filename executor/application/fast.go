package executorapplication

import (
	"strings"
	"sync"

	"github.com/acarl005/stripansi"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type FastExecutor struct {
	shell   shelldomain.Shell
	locker  sync.Locker
	wrapper wrapperdomain.CmdWrapper
	reader  executordomain.Reader
}

func NewFastExecutor(shell shelldomain.Shell, locker sync.Locker, wrapper wrapperdomain.CmdWrapper) *FastExecutor {
	return &FastExecutor{
		shell:   shell,
		locker:  locker,
		wrapper: wrapper,
		reader:  NewReader(shell),
	}
}

func (e *FastExecutor) Run(command string) (string, int, error) {
	e.locker.Lock()
	defer e.locker.Unlock()

	wrappedCommand := e.wrapper.WrapCmd(command)
	_, err := e.shell.Write([]byte(wrappedCommand))
	if err != nil {
		return err.Error(), -2, err
	}

	started := false
	code := -10

	output, err := e.reader.Read(func(s string) (string, bool) {
		result := e.wrapper.UnwrapCmd(s, started)
		started = result.Started
		code = result.Code
		return result.Output, !result.IsRunning
	})

	return strings.TrimSpace(output), code, err
}

func (e *FastExecutor) RunAndClean(command string) (string, int, error) {
	output, code, err := e.Run(command)
	cleanText := stripansi.Strip(output)
	return cleanText, code, err
}
