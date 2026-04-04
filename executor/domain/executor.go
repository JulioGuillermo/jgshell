package executordomain

type Executor interface {
	Run(command string) (*Cmd, error)
	IsRunning() bool
}
