package executorapplication

import (
	"sync"
	"time"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	toolsdomain "github.com/julioguillermo/jgshell/tools/domain"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type Executor struct {
	shell     shelldomain.Shell
	reader    executordomain.Reader
	uuid      toolsdomain.UUIDGenerator
	wrapper   wrapperdomain.CmdWrapper
	locker    sync.Locker
	cond      *sync.Cond
	isRunning bool
	cmd       *executordomain.Cmd
}

func NewExecutor(
	shell shelldomain.Shell,
	locker sync.Locker,
	wrapper wrapperdomain.CmdWrapper,
	uuid toolsdomain.UUIDGenerator,
) *Executor {
	e := &Executor{
		shell:   shell,
		locker:  locker,
		wrapper: wrapper,
		uuid:    uuid,
		reader:  NewReader(shell),
		cond:    sync.NewCond(locker),
	}
	e.startReader()
	return e
}

func (e *Executor) IsRunning() bool {
	return e.isRunning
}

func (e *Executor) Run(command string) (*executordomain.Cmd, error) {
	if e.isRunning {
		_, err := e.shell.Write([]byte(command))
		return nil, err
	}
	return e.runNewCmd(command)
}

func (e *Executor) runNewCmd(command string) (*executordomain.Cmd, error) {
	e.locker.Lock()
	defer e.locker.Unlock()
	defer e.cond.Signal()

	start := time.Now()
	e.cmd = &executordomain.Cmd{
		UUID:     e.uuid.Generate(),
		Cmd:      command,
		Start:    &start,
		ExitCode: -10,
	}

	e.isRunning = true

	cmd := e.wrapper.WrapCmd(command)
	_, err := e.shell.Write([]byte(cmd))
	if err != nil {
		e.isRunning = false

		end := time.Now()
		e.cmd.End = &end
		e.cmd.Error = err
		e.cmd.ExitCode = -2

		return e.cmd, err
	}

	return e.cmd, nil
}

func (e *Executor) startReader() {
	go func() {
		for {
			_, err := e.reader.ReadPrecond(
				e.cond.L,
				e.preCond,
				e.processOutput,
			)

			if err != nil && e.cmd != nil {
				e.cmd.Error = err
			}
		}
	}()
}

func (e *Executor) preCond(string) bool {
	for !e.isRunning {
		e.cond.Wait()
	}
	return false
}

func (e *Executor) processOutput(output string) (string, bool) {
	result := e.wrapper.UnwrapCmd(output, e.cmd.Started)

	e.cmd.Output = result.Output
	e.cmd.ExitCode = result.Code

	if result.Started {
		e.cmd.Started = true
	}

	if e.cmd.PWD == "" {
		e.cmd.PWD = result.Pwd
	}
	if e.cmd.USER == "" {
		e.cmd.USER = result.User
	}

	if result.IsRunning {
		return output, false
	}

	e.cmd.End = result.EndTime
	e.isRunning = false

	return output, true
}
