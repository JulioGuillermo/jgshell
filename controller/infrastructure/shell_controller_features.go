package controllerinfrastructure

import statusdomain "github.com/julioguillermo/jgshell/status/domain"

func (ctl *ShellController) GetStatus() (*statusdomain.Status, error) {
	return ctl.statusLoader.Load()
}

func (ctl *ShellController) GetAutocomplete(line string, cursor int) ([]string, error) {
	return ctl.autocomplete.GetAutocomplete(line, cursor)
}
