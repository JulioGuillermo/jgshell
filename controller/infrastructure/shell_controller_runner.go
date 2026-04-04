package controllerinfrastructure

import (
	"errors"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
)

func (ctl *ShellController) IsRunning() bool {
	if ctl.shellExecutor == nil {
		return false
	}
	return ctl.shellExecutor.IsRunning()
}

func (ctl *ShellController) Run(command string) error {
	if ctl.shellExecutor == nil {
		return errors.New("Fail to run command: shell executor is not initialized")
	}
	cmd, err := ctl.shellExecutor.Run(command)
	if err != nil {
		return err
	}
	ctl.history.PushCmd(cmd)
	return nil
}

func (ctl *ShellController) GetHistory() []*executordomain.Cmd {
	return ctl.history.GetCmd()
}
