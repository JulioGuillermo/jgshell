package autocompleteinfrastructure

type autocompleteConfig struct {
	cmd    string
	script string
}

var autocompleteConfs = map[string]autocompleteConfig{
	"bash":       {"bash", "autocomplete/autocomplete_bash.sh"},
	"zsh":        {"zsh", "autocomplete/autocomplete_zsh.sh"},
	"fish":       {"fish", "autocomplete/autocomplete_fish.fish"},
	"powershell": {"pwsh", "autocomplete/autocomplete_powershell.ps1"},
	"nushell":    {"nu", "autocomplete/autocomplete_nushell.nu"},
	"sh":         {"sh", "autocomplete/autocomplete_sh.sh"},
}

func (a *Autocomplete) getConfigForShell(sh string) autocompleteConfig {
	conf, ok := autocompleteConfs[sh]
	if !ok {
		conf = autocompleteConfs["bash"]
	}
	return conf
}
