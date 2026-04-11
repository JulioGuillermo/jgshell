package executorapplication

import (
	"fmt"
	"strings"
	"sync"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	routerdomain "github.com/julioguillermo/jgshell/router/domain"
	toolsdomain "github.com/julioguillermo/jgshell/tools/domain"
	wrapperdomain "github.com/julioguillermo/jgshell/wrapper/domain"
)

type SimpleExecutor struct {
	router        routerdomain.Router
	locker        sync.Locker
	uuidGenerator toolsdomain.UUIDGenerator
}

func NewSimpleExecutor(router routerdomain.Router, uuidGenerator toolsdomain.UUIDGenerator) *SimpleExecutor {
	return &SimpleExecutor{
		locker:        &sync.Mutex{},
		router:        router,
		uuidGenerator: uuidGenerator,
	}
}

func (s *SimpleExecutor) Run(command string) (string, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	uuid := s.uuidGenerator.Generate()

	s.router.ClearQueue(
		executordomain.SimpleQueue,
	)
	_, err := fmt.Fprintf(
		s.router,
		`%s "%s" "$(%s)" "%s"
`,
		wrapperdomain.SimpleWrapper,
		uuid, command, uuid,
	)
	if err != nil {
		return "", err
	}

	element, err := s.router.ReadFrom(executordomain.SimpleQueue)
	if err != nil {
		return "", err
	}

	output := element.FinalString()

	Start := fmt.Sprintf(`<<<JGSHELL_START;%s>>>`, uuid)
	End := fmt.Sprintf(`<<<JGSHELL_END;%s>>>`, uuid)

	idx := strings.Index(output, Start)
	if idx != -1 {
		output = output[idx+len(Start):]
	}
	idx = strings.Index(output, End)
	if idx != -1 {
		output = output[:idx]
	}

	return strings.TrimSpace(output), nil
}
