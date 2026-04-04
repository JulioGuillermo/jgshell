package controllerinfrastructure

import statusdomain "github.com/julioguillermo/jgshell/status/domain"

func (ctl *ShellController) GetStatus() (*statusdomain.Status, error) {
	return ctl.statusLoader.Load()
}

func (ctl *ShellController) GetAutocomplete(line string, cursor int) ([]string, error) {
	return ctl.autocomplete.GetAutocomplete(line, cursor)
}

func (ctl *ShellController) GetCmdHistory() []string {
	return ctl.persistencer.Get()
}

func (ctl *ShellController) Filter(start string) []string {
	return ctl.persistencer.Filter(start)
}

func (ctl *ShellController) FilterLast(start string) string {
	return ctl.persistencer.FilterLast(start)
}
