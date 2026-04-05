package persistenceapplication

import (
	"slices"
	"strings"

	persistencedomain "github.com/julioguillermo/jgshell/persistence/domain"
)

const maxHistory = 5000

type PersistenceController struct {
	cmds        []string
	persistence persistencedomain.Persistence
}

func NewPersistenceController(persistence persistencedomain.Persistence) (*PersistenceController, error) {
	history, err := persistence.LoadHistory()
	if err != nil {
		return nil, err
	}
	return &PersistenceController{
		cmds:        history,
		persistence: persistence,
	}, nil
}

func (p *PersistenceController) Push(cmd string) error {
	p.cmds = slices.DeleteFunc(p.cmds, func(c string) bool {
		return c == cmd
	})

	p.cmds = append(p.cmds, cmd)
	if len(p.cmds) > maxHistory {
		p.cmds = p.cmds[len(p.cmds)-maxHistory:]
	}
	return p.persistence.SaveHistory(p.cmds)
}

func (p *PersistenceController) Get() []string {
	cp := make([]string, len(p.cmds))
	copy(cp, p.cmds)
	return cp
}

func (p *PersistenceController) Filter(start string) []string {
	return slices.DeleteFunc(p.Get(), func(c string) bool {
		return !strings.HasPrefix(c, start)
	})
}

func (p *PersistenceController) FilterLast(start string) string {
	if len(p.cmds) == 0 {
		return ""
	}
	for i := len(p.cmds) - 1; i >= 0; i-- {
		if strings.HasPrefix(p.cmds[i], start) {
			return p.cmds[i]
		}
	}
	return ""
}
