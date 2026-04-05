package executorapplication

import (
	"sync"

	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
)

type Reader struct {
	shell    shelldomain.Shell
	cleanner *Cleanner
}

func NewReader(shell shelldomain.Shell) *Reader {
	return &Reader{
		shell:    shell,
		cleanner: NewCleanner(),
	}
}

func (r *Reader) ReadPrecond(locker sync.Locker, pre func(string) bool, post func(string) (string, bool)) (string, error) {
	var stop bool
	var err error
	var n int
	output := ""
	buf := make([]byte, 1024)

	for !func() bool {
		locker.Lock()
		defer locker.Unlock()
		if pre(output) {
			return true
		}

		locker.Unlock()
		n, err = r.shell.Read(buf)
		locker.Lock()
		if err != nil {
			return true
		}
		if n <= 0 {
			return false
		}

		output, stop = post(r.Clear(output + string(buf[:n])))
		return stop
	}() {
	}

	return output, err
}

func (r *Reader) Read(f func(string) (string, bool)) (string, error) {
	var stop bool
	output := ""
	buf := make([]byte, 1024)

	// re := regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

	for {
		n, err := r.shell.Read(buf)
		if err != nil {
			return output, err
		}
		if n <= 0 {
			continue
		}
		output, stop = f(r.Clear(output + string(buf[:n])))
		// fmt.Println(re.ReplaceAllString(output, ""))
		if stop {
			break
		}
	}

	return output, nil
}

func (r *Reader) Clear(str string) string {
	return str
	// if r.cleanner == nil {
	// 	r.cleanner = NewCleanner()
	// }
	// return r.cleanner.Clear(str)
}
