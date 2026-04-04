package executordomain

type SimpleExecutor interface {
	Run(string) (string, error)
}
