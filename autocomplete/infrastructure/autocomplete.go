package autocompleteinfrastructure

import (
	"strings"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
)

type Autocomplete struct {
	shellDetector shelldetectordomain.ShellDetector
	executor      executordomain.FastExecutor
}

func NewAutocomplete(
	shellDetector shelldetectordomain.ShellDetector,
	executor executordomain.FastExecutor,
) *Autocomplete {
	return &Autocomplete{
		shellDetector: shellDetector,
		executor:      executor,
	}
}

func (a *Autocomplete) GetAutocomplete(line string, cursor int) ([]string, error) {
	sh, err := a.shellDetector.DetectShell()
	if err != nil {
		return nil, err
	}
	conf := a.getConfigForShell(sh)
	script, err := a.getScript(conf.script)
	if err != nil {
		return nil, err
	}

	cmd := a.buildCmd(cursor, line, sh, conf.cmd, script)

	output, code, err := a.executor.RunAndClean(cmd)
	if err != nil {
		return nil, err
	}
	if code != 0 {
		return []string{}, nil
	}

	var completions []string
	for l := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			completions = append(completions, l)
		}
	}
	return completions, nil
}
