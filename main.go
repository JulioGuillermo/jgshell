package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/app"
	controllerdomain "github.com/julioguillermo/jgshell/controller/domain"
	controllerinfrastructure "github.com/julioguillermo/jgshell/controller/infrastructure"
	statusdomain "github.com/julioguillermo/jgshell/status/domain"
	syntaxinfrastructure "github.com/julioguillermo/jgshell/syntax/infrastruct"
)

func main() {
	cmd := ""
	status := false
	shell := false
	for i, a := range os.Args {
		if i == 0 {
			continue
		}
		if i == 1 {
			switch a {
			case "--status":
				status = true
				continue
			case "--shell":
				shell = true
				continue
			}
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

	if shell {
		shell, err := ctl.GetShell()
		if err != nil {
			fmt.Printf("Fail to get shell: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Shell: %s\n", shell)
		return
	}

	err = ctl.WrapShell()
	if err != nil {
		fmt.Printf("Fail to wrap shell: %v", err)
		os.Exit(1)
	}

	if status {
		shell, err := ctl.GetShell()
		if err != nil {
			fmt.Printf("Fail to get shell: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Shell: %s\n", shell)
		out, exit, err := ctl.ShellFastExecutor.Run("printf 'Hello World'")
		if err != nil {
			fmt.Printf("Fail to run command: %v", err)
			os.Exit(1)
		}
		fmt.Println(exit)
		fmt.Println(out)

		status, err := ctl.GetStatus()
		if err != nil {
			fmt.Printf("Fail to get status: %v", err)
			os.Exit(1)
		}
		if status.Git == nil {
			status.Git = &statusdomain.Git{}
		}
		fmt.Printf(
			"OS: %s\nShell: %s\nUser: %s\nDir: %s\n"+
				"GIT Branch: %s\nGIT Branch Remote: %s\nGIT Ahead: %d\nGIT Behind: %d\n"+
				"GIT Untracked: %d\nGIT Modified: %d\nGIT Staged: %d\nGIT Deleted: %d\nGIT Conflicts: %d\n",
			status.OS,
			status.Shell,
			status.User,
			status.Dir,
			status.Git.BranchLocal,
			status.Git.BranchRemote,
			status.Git.Ahead,
			status.Git.Behind,
			status.Git.Untracked,
			status.Git.Modified,
			status.Git.Staged,
			status.Git.Deleted,
			status.Git.Conflicts,
		)
		return
	}

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
