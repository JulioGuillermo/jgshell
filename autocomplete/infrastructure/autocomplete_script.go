package autocompleteinfrastructure

import "github.com/julioguillermo/jgshell/scripts"

func (a *Autocomplete) getScript(path string) (string, error) {
	bytes, err := scripts.AutoCompleteScript.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
