package executordomain

type FastExecutor interface {
	Run(command string) (string, int, error)
	RunAndClean(command string) (string, int, error)
}
