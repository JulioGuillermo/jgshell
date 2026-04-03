package statedomain

type State interface {
	Send(string) error
	Write([]byte) error
	GetHistory() []Cmd

	SetSize(width, height int)
	IsRunning() bool

	GetAutoComplete(line string, cursor int) []string
	GetStatus() Status

	OnClose(f func(s State))
}
