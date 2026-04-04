package executordomain

type Executor interface {
	Run(command string) (*Cmd, error)
	IsRunning() bool
	Stop()
	StopWith(code int, msg string)
}
