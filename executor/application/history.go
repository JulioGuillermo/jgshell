package executorapplication

import (
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
)

type History struct {
	cmds []*executordomain.Cmd
}

func NewHistory() *History {
	h := &History{}
	h.Clear()
	return h
}

func (h *History) PushCmd(cmd *executordomain.Cmd) {
	h.cmds = append(h.cmds, cmd)
}

func (h *History) GetCmd() []*executordomain.Cmd {
	return h.cmds
}

func (h *History) LastCmd() *executordomain.Cmd {
	if len(h.cmds) == 0 {
		return nil
	}
	return h.cmds[len(h.cmds)-1]
}

func (h *History) Clear() {
	h.cmds = []*executordomain.Cmd{}
}
