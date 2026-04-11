package routerdomain

import "io"

type Router interface {
	io.Writer
	WriteBytes([]byte) error
	WriteString(string) error
	ReadFrom(string) (Element, error)
	ClearQueue(name string)
	Reset()
}
