package executordomain

type History interface {
	PushCmd(cmd *Cmd)
	GetCmd() []*Cmd
	Clear()
}
