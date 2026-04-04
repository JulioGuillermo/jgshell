package persistenceinfrastructure

import (
	"os"
	"path"
	"strings"
)

type Persistence struct {
}

func NewPersistence() *Persistence {
	return &Persistence{}
}

func (p *Persistence) getHistoryFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := path.Join(home, ".jgshell")
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return "", err
	}
	return path.Join(dir, "history"), nil
}

func (p *Persistence) LoadHistory() ([]string, error) {
	historyFilePath, err := p.getHistoryFilePath()
	if err != nil {
		return nil, err
	}
	history, err := os.ReadFile(historyFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	return strings.Split(string(history), "\n"), nil
}

func (p *Persistence) SaveHistory(history []string) error {
	historyFilePath, err := p.getHistoryFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(
		historyFilePath,
		[]byte(strings.Join(history, "\n")),
		0777,
	)
}
