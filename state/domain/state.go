package statedomain

type State interface {
	Send(string) error
	GetHistory() []Cmd

	SetSize(width, height int)
	IsRunning() bool

	GetAutoComplete() []string
	GetStatus()

	OnClose(f func(s State))
}
