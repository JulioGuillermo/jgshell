package statedomain

type FastCmd interface {
	GetShell() string
	FastCmd(string) (string, int)
	FastCmdClean(string) (string, int)
}
