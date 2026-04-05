package controllerinfrastructure

import (
	"errors"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
)

func (ctl *ShellController) IsRunning() bool {
	if ctl.ShellExecutor == nil {
		return false
	}
	return ctl.ShellExecutor.IsRunning()
}

func (ctl *ShellController) Run(command string) error {
	if ctl.ShellExecutor == nil {
		return errors.New("Fail to run command: shell executor is not initialized")
	}
	cmd, err := ctl.ShellExecutor.Run(command)
	if err != nil {
		return err
	}
	if cmd != nil {
		ctl.History.PushCmd(cmd)
		err = ctl.Persistencer.Push(cmd.Cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctl *ShellController) GetHistory() []*executordomain.Cmd {
	return ctl.History.GetCmd()
}

func (ctl *ShellController) LastCmd() *executordomain.Cmd {
	return ctl.History.LastCmd()
}
