package wrapperdomain

type CmdWrapper interface {
	WrapCmd(cmd string) string
	UnwrapCmd(output string, started bool) *CmdUnwrapResult
}
