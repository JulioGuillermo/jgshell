package persistenceinfrastructure

import (
	persistenceapplication "github.com/julioguillermo/jgshell/persistence/application"
	persistencedomain "github.com/julioguillermo/jgshell/persistence/domain"
)

func NewPersistenceCtl() (persistencedomain.PersistenceCtl, error) {
	persistence := NewPersistence()
	return persistenceapplication.NewPersistenceController(persistence)
}
