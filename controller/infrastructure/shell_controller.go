package controllerinfrastructure

import (
	"sync"

	executorapplication "github.com/julioguillermo/jgshell/executor/application"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	shellinfrastructure "github.com/julioguillermo/jgshell/shell/infrastructure"
	shelldetectorapplication "github.com/julioguillermo/jgshell/shelldetector/application"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
	toolsapplication "github.com/julioguillermo/jgshell/tools/application"
	toolsdomain "github.com/julioguillermo/jgshell/tools/domain"
	toolsinfrastructure "github.com/julioguillermo/jgshell/tools/infrastructure"
	wrapperapplication "github.com/julioguillermo/jgshell/wrapper/application"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
	wrapperinfrastructure "github.com/julioguillermo/jgshell/wrapper/infrastructure"
)

type ShellController struct {
	locker sync.Locker

	uuidGenerator  toolsdomain.UUIDGenerator
	outputCleanner toolsdomain.OutputCleaner

	history executordomain.History

	shell               shelldomain.FullShell
	shellDetector       shelldetectordomain.ShellDetector
	shellWrapper        wrapperdomain.ShellWrapper
	shellCmdWrapper     wrapperdomain.CmdWrapper
	shellSimpleExecutor executordomain.SimpleExecutor
	shellFastExecutor   executordomain.FastExecutor
	shellExecutor       executordomain.Executor
}

func NewShellController(cmd string) (*ShellController, error) {
	ctl := &ShellController{
		locker:         &sync.Mutex{},
		uuidGenerator:  toolsinfrastructure.NewUUIDGenerator(),
		outputCleanner: toolsapplication.NewOutputCleaner(),
		history:        executorapplication.NewHistory(),
	}

	err := ctl.initShell(cmd)
	if err != nil {
		return nil, err
	}

	err = ctl.initExecutors()
	if err != nil {
		return nil, err
	}

	return ctl, nil
}

func (ctl *ShellController) initShell(cmd string) error {
	sh, err := shellinfrastructure.NewShellConnector(cmd)
	if err != nil {
		return err
	}
	ctl.shell = sh
	return nil
}

func (ctl *ShellController) initExecutors() error {
	ctl.shellSimpleExecutor = executorapplication.NewSimpleExecutor(ctl.shell, ctl.locker, ctl.uuidGenerator)
	ctl.shellDetector = shelldetectorapplication.NewShellDetector(ctl.shellSimpleExecutor)
	ctl.shellWrapper = wrapperinfrastructure.NewShellWrapper(ctl.shell, ctl.shellDetector)
	ctl.shellCmdWrapper = wrapperapplication.NewCmdWrapper()
	ctl.shellFastExecutor = executorapplication.NewFastExecutor(ctl.shell, ctl.locker, ctl.shellCmdWrapper)
	ctl.shellExecutor = executorapplication.NewExecutor(ctl.shell, ctl.locker, ctl.shellCmdWrapper, ctl.uuidGenerator)

	return nil
}
