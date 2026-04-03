package stateapplication

import (
	"time"

	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func (s *State) Send(message string) error {
	return s.Write([]byte(message))
}

func (s *State) Write(message []byte) error {
	if s.isRunning {
		_, err := s.shell.Write(message)
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	defer s.cond.Signal()

	uuid := GetUUID()
	s.isRunning = true
	cmd := WrapCmd(string(message), uuid)

	start := time.Now()
	s.history = append(s.history, statedomain.Cmd{
		UUID:  uuid,
		Cmd:   string(message),
		Start: &start,
	})

	_, err := s.shell.Write([]byte(cmd))
	if err != nil {
		s.isRunning = false
		return err
	}

	return nil
}

func (s *State) startReader() {
	s.shell.Write([]byte(GetPS1()))
	go func() {
		for {
			s.readOutput()
		}
	}()
}

func (s *State) readOutput() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for !s.isRunning {
		s.cond.Wait()
	}

	buffer := make([]byte, 1024)
	lastCmd := s.lastCmd()
	n, err := s.shell.Read(buffer)

	if lastCmd == nil {
		return
	}
	if err != nil {
		lastCmd.Error = err
		return
	}

	if n > 0 {
		result := CleanOutput(lastCmd.Output+string(buffer[:n]), lastCmd.UUID)
		lastCmd.Output = result.Output
		lastCmd.ExitCode = result.Code
		if result.Started {
			lastCmd.Started = true
		}
		if lastCmd.PWD == "" {
			lastCmd.PWD = result.Pwd
		}
		if lastCmd.USER == "" {
			lastCmd.USER = result.Username
		}

		if result.IsRunning {
			return
		}
		lastCmd.End = result.EndTime
		s.isRunning = false
	}
}
