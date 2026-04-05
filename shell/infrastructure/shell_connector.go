package shellinfrastructure

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
)

type ShellConnector struct {
	cmdName string
	cmdArgs []string
	cmd     *exec.Cmd
	ptyFile *os.File
}

func NewShellConnector(cmd string) (shelldomain.FullShell, error) {
	sh := &ShellConnector{}
	if err := sh.start(cmd); err != nil {
		return nil, err
	}

	return sh, nil
}

func (s *ShellConnector) start(cmd string) error {
	args := strings.Fields(cmd)
	s.cmdName = args[0]
	s.cmdArgs = args[1:]

	s.cmd = exec.Command(s.cmdName, s.cmdArgs...)

	if err := s.initEnv(); err != nil {
		return err
	}

	if err := s.setDir(); err != nil {
		return err
	}

	if err := s.initCmd(); err != nil {
		return err
	}

	return nil
}

func (s *ShellConnector) initEnv() error {
	if s.cmd.Env == nil {
		s.cmd.Env = []string{}
	}

	s.cmd.Env = append(s.cmd.Env, os.Environ()...)
	s.cmd.Env = append(
		s.cmd.Env,
		// "TERM=dumb",
		"TERM=xterm",
		// "TERM=xterm-256color",
		"GIT_TERMINAL_PROMPT=1",
		"GIT_CPT_FORBID_DECORATION=1",
	)

	return nil
}

func (s *ShellConnector) setDir() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	s.cmd.Dir = dir

	return nil
}

func (s *ShellConnector) initCmd() error {
	ptyFile, err := pty.Start(s.cmd)
	if err != nil {
		return err
	}
	err = configPty(ptyFile.Fd())
	if err != nil {
		return err
	}
	// fd := s.ptyFile.Fd()
	// _, err = term.MakeRaw(fd)
	// if err != nil {
	// 	return err
	// }
	s.ptyFile = ptyFile
	return s.SetSize(24, 80)
}

func (s *ShellConnector) SetSize(r, c uint16) error {
	return pty.Setsize(s.ptyFile, &pty.Winsize{
		Rows: r,
		Cols: c,
	})
}

func (s *ShellConnector) Write(p []byte) (int, error) {
	if s.ptyFile == nil {
		return 0, io.ErrClosedPipe
	}
	return s.ptyFile.Write(p)
}

func (s *ShellConnector) Read(p []byte) (int, error) {
	if s.ptyFile == nil {
		return 0, io.ErrClosedPipe
	}
	n, err := s.ptyFile.Read(p)
	if n > 0 {
		data := string(p[:n])
		if strings.Contains(data, "\x1b[6n") {
			s.ptyFile.Write([]byte("\x1b[2;2R"))
		}
	}
	return n, err
	// return s.ptyFile.Read(p)
}

func (s *ShellConnector) Close() error {
	if s.ptyFile != nil {
		s.ptyFile.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}

func (s *ShellConnector) OnClose(f func(shelldomain.FullShell)) {
	go func() {
		s.cmd.Wait()
		f(s)
	}()
}
