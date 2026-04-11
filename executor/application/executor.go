package executorapplication

import (
	"sync"
	"time"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	routerdomain "github.com/julioguillermo/jgshell/router/domain"
	shelldetectordomain "github.com/julioguillermo/jgshell/shelldetector/domain"
	toolsdomain "github.com/julioguillermo/jgshell/tools/domain"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type Executor struct {
	router        routerdomain.Router
	shellDetector shelldetectordomain.ShellDetector
	uuid          toolsdomain.UUIDGenerator
	wrapper       wrapperdomain.CmdWrapper
	locker        sync.Locker
	cond          *sync.Cond
	isRunning     bool
	cmd           *executordomain.Cmd
}

func NewExecutor(
	router routerdomain.Router,
	shellDetector shelldetectordomain.ShellDetector,
	wrapper wrapperdomain.CmdWrapper,
	uuid toolsdomain.UUIDGenerator,
) *Executor {
	locker := &sync.Mutex{}
	e := &Executor{
		router:        router,
		shellDetector: shellDetector,
		locker:        locker,
		wrapper:       wrapper,
		uuid:          uuid,
		cond:          sync.NewCond(locker),
	}
	e.startReader()
	return e
}

func (e *Executor) IsRunning() bool {
	return e.isRunning
}

func (e *Executor) StopWith(code int, msg string) {
	e.isRunning = false

	if e.cmd == nil || !e.cmd.IsRunning() {
		return
	}

	end := time.Now()
	e.cmd.End = &end
	e.cmd.ExitCode = code
	e.cmd.Output += "\n" + msg
}

func (e *Executor) Stop() {
	e.StopWith(-11, "Stopped here")
}

func (e *Executor) Run(command string) (*executordomain.Cmd, error) {
	if e.isRunning {
		err := e.router.WriteString(command)
		return nil, err
	}
	return e.runNewCmd(command)
}

func (e *Executor) runNewCmd(command string) (*executordomain.Cmd, error) {
	e.locker.Lock()
	defer e.locker.Unlock()
	defer e.cond.Signal()

	sh, err := e.shellDetector.DetectShell()
	if err != nil {
		return nil, err
	}

	start := time.Now()
	e.cmd = &executordomain.Cmd{
		SH:       sh,
		UUID:     e.uuid.Generate(),
		Cmd:      command,
		Start:    &start,
		ExitCode: -10,
	}

	e.isRunning = true

	e.router.ClearQueue(executordomain.DefaultQueue)
	cmd := e.wrapper.WrapCmd(sh, command)
	err = e.router.WriteString(cmd)
	if err != nil {
		e.cmdError(err)
		return e.cmd, err
	}

	return e.cmd, nil
}

func (e *Executor) cmdError(err error) {
	if e.cmd == nil {
		return
	}

	e.isRunning = false

	end := time.Now()
	e.cmd.End = &end
	e.cmd.Error = err
	e.cmd.ExitCode = -2
}

func (e *Executor) startReader() {
	go func() {
		for {
			e.read()
		}
	}()
}

func (e *Executor) read() {
	e.locker.Lock()
	defer e.locker.Unlock()
	for !e.isRunning {
		e.cond.Wait()
	}

	element, err := e.router.ReadFrom(executordomain.DefaultQueue)
	if err != nil {
		e.cmdError(err)
		return
	}

	for e.isRunning && !element.IsEnded() {
		e.processOutput(element.String())
		time.Sleep(time.Millisecond)
	}

	e.processOutput(element.String())

	end := time.Now()
	e.cmd.End = &end
	e.isRunning = false
}

func (e *Executor) processOutput(output string) {
	if e.cmd == nil {
		return
	}

	result := e.wrapper.UnwrapCmd(output, false)

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
}
