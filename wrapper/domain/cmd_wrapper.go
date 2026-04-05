package wrapperdomain

type CmdWrapper interface {
	WrapCmd(sh, cmd string) string
	UnwrapCmd(output string, started bool) *CmdUnwrapResult
}
