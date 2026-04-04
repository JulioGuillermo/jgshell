package shelldomain

import "io"

type Shell interface {
	io.Writer
	io.Reader
}

type FullShell interface {
	Shell
	SetSize(r, c uint16) error
	io.Closer
	OnClose(f func(FullShell))
}
