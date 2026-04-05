package controllerinfrastructure

import (
	"errors"

	controllerdomain "github.com/julioguillermo/jgshell/controller/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
)

func (ctl *ShellController) GetShell() (string, error) {
	if ctl.ShellDetector == nil {
		return "", errors.New("Shell detector not initialized")
	}
	return ctl.ShellDetector.DetectShell()
}

func (ctl *ShellController) WrapShell() error {
	if ctl.ShellWrapper == nil {
		return errors.New("Fail to wrap uninitialized shell")
	}
	ctl.ShellExecutor.StopWith(-12, "Stop and wrap shell")
	return ctl.ShellWrapper.WrapShell()
}

func (ctl *ShellController) SetSize(width, height int) error {
	if ctl.Shell == nil {
		return nil
	}
	return ctl.Shell.SetSize(uint16(height), uint16(width))
}

func (ctl *ShellController) OnClose(f func(controllerdomain.ShellController)) {
	if ctl.Shell == nil {
		return
	}
	ctl.Shell.OnClose(func(fs shelldomain.FullShell) { f(ctl) })
}

func (ctl *ShellController) Close() error {
	if ctl.Shell == nil {
		return nil
	}
	return ctl.Shell.Close()
}
