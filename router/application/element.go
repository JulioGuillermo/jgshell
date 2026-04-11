package routerapplication

import "context"

type Element struct {
	data   string
	ctx    context.Context
	cancel context.CancelFunc
	ended  bool
}

func NewElement() *Element {
	ctx, cancel := context.WithCancel(context.Background())
	return &Element{
		data:   "",
		ctx:    ctx,
		cancel: cancel,
	}
}

func (e *Element) AppendData(data string) {
	e.data += data
}

func (e *Element) String() string {
	return e.data
}

func (e *Element) FinalString() string {
	<-e.ctx.Done()
	return e.data
}

func (e *Element) IsEnded() bool {
	return e.ended
}

func (e *Element) Close() {
	e.cancel()
	e.ended = true
}
