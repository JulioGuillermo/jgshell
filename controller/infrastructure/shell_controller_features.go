package controllerinfrastructure

import statusdomain "github.com/julioguillermo/jgshell/status/domain"

func (ctl *ShellController) GetStatus() (*statusdomain.Status, error) {
	return ctl.StatusLoader.Load()
}

func (ctl *ShellController) GetAutocomplete(line string, cursor int) ([]string, error) {
	return ctl.Autocomplete.GetAutocomplete(line, cursor)
}

func (ctl *ShellController) GetCmdHistory() []string {
	return ctl.Persistencer.Get()
}

func (ctl *ShellController) Filter(start string) []string {
	return ctl.Persistencer.Filter(start)
}

func (ctl *ShellController) FilterLast(start string) string {
	return ctl.Persistencer.FilterLast(start)
}
