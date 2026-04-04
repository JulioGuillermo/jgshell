package executorapplication

import shelldomain "github.com/julioguillermo/jgshell/shell/domain"

type Reader struct {
	shell shelldomain.Shell
}

func NewReader(shell shelldomain.Shell) *Reader {
	return &Reader{shell: shell}
}

func (r *Reader) Read(f func(string) (string, bool)) (string, error) {
	var stop bool
	output := ""
	buf := make([]byte, 1024)

	for {
		n, err := r.shell.Read(buf)
		if err != nil {
			return output, err
		}
		if n <= 0 {
			continue
		}
		output, stop = f(output + string(buf[:n]))
		if stop {
			break
		}
	}

	return output, nil
}
