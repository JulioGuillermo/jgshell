package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/app"
	controllerdomain "github.com/julioguillermo/jgshell/controller/domain"
	controllerinfrastructure "github.com/julioguillermo/jgshell/controller/infrastructure"
	syntaxinfrastructure "github.com/julioguillermo/jgshell/syntax/infrastruct"
)

func main() {
	cmd := ""
	for i, a := range os.Args {
		if i == 0 {
			continue
		}
		if a != "" {
			if i > 1 {
				cmd += " "
			}
			cmd += a
		}
	}
	if cmd == "" {
		cmd = "bash"
	}

	ctl, err := controllerinfrastructure.NewShellController(cmd)
	if err != nil {
		fmt.Printf("Fail to create controller: %v", err)
		os.Exit(1)
	}
	defer ctl.Close()
	ctl.SetSize(24, 80)
	err = ctl.WrapShell()
	if err != nil {
		fmt.Printf("Fail to wrap shell: %v", err)
		os.Exit(1)
	}
	ctl.GetStatus()

	hl, err := syntaxinfrastructure.NewTSHighlighter()
	if err != nil {
		fmt.Printf("Fail to create highlighter: %v", err)
		os.Exit(1)
	}
	app := app.NewApp(ctl, hl)

	p := tea.NewProgram(app)
	ctl.OnClose(func(sc controllerdomain.ShellController) {
		p.Quit()
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Fail to run the app: %v", err)
		os.Exit(1)
	}
}
