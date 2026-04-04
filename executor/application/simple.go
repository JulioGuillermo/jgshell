package executorapplication

import (
	"fmt"
	"strings"
	"sync"

	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	toolsdomain "github.com/julioguillermo/jgshell/tools/domain"
)

type SimpleExecutor struct {
	shell         shelldomain.Shell
	mu            *sync.Mutex
	uuidGenerator toolsdomain.UUIDGenerator
	reader        executordomain.Reader
}

func NewSimpleExecutor(shell shelldomain.Shell, mu *sync.Mutex, uuidGenerator toolsdomain.UUIDGenerator) *SimpleExecutor {
	return &SimpleExecutor{
		shell:         shell,
		mu:            mu,
		uuidGenerator: uuidGenerator,
		reader:        NewReader(shell),
	}
}

func (s *SimpleExecutor) Run(command string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	uuid := s.uuidGenerator.Generate()

	_, err := fmt.Fprintf(
		s.shell,
		`printf "<<<JGSHELL_START;%%s>>> %%s <<<JGSHELL_END;%%s>>>\n" "%s" "$(%s)" "%s"
`,
		uuid, command, uuid,
	)
	if err != nil {
		return "", err
	}

	Start := fmt.Sprintf(`<<<JGSHELL_START;%s>>>`, uuid)
	End := fmt.Sprintf(`<<<JGSHELL_END;%s>>>`, uuid)

	output, err := s.reader.Read(func(s string) (string, bool) {
		return s, strings.Contains(s, End)
	})
	if err != nil {
		return "", err
	}

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
