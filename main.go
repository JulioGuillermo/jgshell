package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/app"
	shellinfrastructure "github.com/julioguillermo/jgshell/shell/infrastructure"
	stateapplication "github.com/julioguillermo/jgshell/state/application"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
	syntaxinfrastructure "github.com/julioguillermo/jgshell/syntax/infrastruct"
)

func main() {
	cmd := ""
	for i, a := range os.Args {
		if i == 0 {
			continue
		}
		if a != "" {
			cmd += a
		}
	}
	if cmd == "" {
		cmd = "bash"
	}

	shell, err := shellinfrastructure.NewShellConnector(cmd)
	if err != nil {
		fmt.Printf("Fail to start shell: %v", err)
		os.Exit(1)
	}
	defer shell.Close()
	shell.SetSize(24, 80)

	state := stateapplication.NewState(shell)

	hl, err := syntaxinfrastructure.NewTSHighlighter()
	if err != nil {
		fmt.Printf("Fail to create highlighter: %v", err)
		os.Exit(1)
	}
	app := app.NewApp(state, hl)

	p := tea.NewProgram(app)
	state.OnClose(func(s statedomain.State) {
		p.Quit()
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fail to run the app: %v", err)
		os.Exit(1)
	}
}
