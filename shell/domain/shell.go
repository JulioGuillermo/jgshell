package shelldomain

import "io"

type Shell interface {
	io.Writer
	io.Reader
	SetSize(r, c uint16) error
	io.Closer
	OnClose(f func(Shell))
}
