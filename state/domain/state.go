package statedomain

type State interface {
	Send(string) error
	GetHistory() []Cmd

	ShowInput() bool
	GetAutoComplete() []string

	ShowStatusBar() bool
	GetDir() string
	GetTime() string

	OnClose(f func(s State))
}
