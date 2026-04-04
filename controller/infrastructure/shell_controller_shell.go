package controllerinfrastructure

import (
	"errors"

	controllerdomain "github.com/julioguillermo/jgshell/controller/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
)

func (ctl *ShellController) WrapShell() error {
	if ctl.shellWrapper == nil {
		return errors.New("Fail to wrap uninitialized shell")
	}
	return ctl.shellWrapper.WrapShell()
}

func (ctl *ShellController) SetSize(width, height int) error {
	if ctl.shell == nil {
		return nil
	}
	return ctl.shell.SetSize(uint16(height), uint16(width))
}

func (ctl *ShellController) OnClose(f func(controllerdomain.ShellController)) {
	if ctl.shell == nil {
		return
	}
	ctl.shell.OnClose(func(fs shelldomain.FullShell) { f(ctl) })
}

func (ctl *ShellController) Close() error {
	if ctl.shell == nil {
		return nil
	}
	return ctl.shell.Close()
}
