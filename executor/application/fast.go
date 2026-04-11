package executorapplication

import (
	"strings"
	"sync"

	"github.com/acarl005/stripansi"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	routerdomain "github.com/julioguillermo/jgshell/router/domain"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type FastExecutor struct {
	router        routerdomain.Router
	shellDetector shelldetectordomain.ShellDetector
	locker        sync.Locker
	wrapper       wrapperdomain.CmdWrapper
}

func NewFastExecutor(
	router routerdomain.Router,
	shellDetector shelldetectordomain.ShellDetector,
	wrapper wrapperdomain.CmdWrapper,
) *FastExecutor {
	return &FastExecutor{
		locker:        &sync.Mutex{},
		router:        router,
		shellDetector: shellDetector,
		wrapper:       wrapper,
	}
}

func (e *FastExecutor) Run(command string) (string, int, error) {
	e.locker.Lock()
	defer e.locker.Unlock()

	sh, err := e.shellDetector.DetectShell()
	if err != nil {
		return "", -1, err
	}

	e.router.ClearQueue(
		executordomain.FastQueue,
	)
	wrappedCommand := e.wrapper.FastWrapCmd(sh, command)
	err = e.router.WriteString(wrappedCommand)
	if err != nil {
		return err.Error(), -2, err
	}

	element, err := e.router.ReadFrom(executordomain.FastQueue)
	if err != nil {
		return err.Error(), -2, err
	}

	output := element.FinalString()
	result := e.wrapper.FastUnwrapCmd(output, false)
	return strings.TrimSpace(result.Output), result.Code, nil
}

func (e *FastExecutor) RunAndClean(command string) (string, int, error) {
	output, code, err := e.Run(command)
	cleanText := stripansi.Strip(output)
	return cleanText, code, err
}
