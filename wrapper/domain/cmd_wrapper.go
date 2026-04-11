package wrapperdomain

type CmdWrapper interface {
	WrapCmd(sh, cmd string) string
	FastWrapCmd(sh, cmd string) string
	UnwrapCmd(output string, started bool) *CmdUnwrapResult
	FastUnwrapCmd(output string, started bool) *CmdUnwrapResult
}
