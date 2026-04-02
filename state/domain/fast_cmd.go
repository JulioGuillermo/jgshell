package statedomain

type FastCmd interface {
	FastCmd(string) (string, int)
	FastCmdClean(string) (string, int)
}
