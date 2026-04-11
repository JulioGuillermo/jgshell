package routerapplication

import (
	"errors"
	"io"
	"sync"

	routerdomain "github.com/julioguillermo/jgshell/router/domain"
)

type Router struct {
	locker     sync.Locker
	readLocker sync.Locker
	reader     io.ReadWriter
	queues     map[string]*Queue

	buffer []byte
	readed string

	queue         string
	queueStartIdx int
	queueEndIdx   int
}

func NewRouter(reader io.ReadWriter, queues ...*Queue) (*Router, error) {
	queueMap := make(map[string]*Queue)
	for _, q := range queues {
		queueMap[q.Name()] = q
	}

	r := &Router{
		locker:     &sync.Mutex{},
		readLocker: &sync.Mutex{},
		reader:     reader,
		queues:     queueMap,
	}
	go r.startReader()
	return r, nil
}

func (r *Router) Write(data []byte) (int, error) {
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.reader.Write(data)
}

func (r *Router) WriteBytes(data []byte) error {
	_, err := r.Write(data)
	return err
}

func (r *Router) WriteString(p string) error {
	return r.WriteBytes([]byte(p))
}

func (r *Router) ReadFrom(name string) (routerdomain.Element, error) {
	queue := r.queues[name]
	if queue == nil {
		return nil, errors.New("Not queue found")
	}
	return queue.Pop(), nil
}

func (r *Router) ClearQueue(name string) {
	queue := r.queues[name]
	if queue == nil {
		return
	}
	queue.Clear()
}

func (r *Router) Reset() {
	r.readLocker.Lock()
	defer r.readLocker.Unlock()

	q := r.queues[r.queue]
	q.End()
	r.readed = ""
	r.queueStartIdx = -1
	r.queueEndIdx = -1
}

func (r *Router) startReader() {
	r.buffer = make([]byte, 1024)
	r.queueStartIdx = -1
	r.queueEndIdx = -1

	for {
		r.read()
	}
}

func (r *Router) read() {
	n, err := r.reader.Read(r.buffer)
	if err != nil {
		return
	}
	if n == 0 {
		return
	}

	r.readLocker.Lock()
	defer r.readLocker.Unlock()

	r.readed += string(r.buffer[:n])

	if r.queueStartIdx == -1 {
		if !r.startQueue() {
			return
		}
		r.readed = r.readed[r.queueStartIdx:]
		r.queues[r.queue].Start()
	}

	if !r.endQueue() {
		r.queues[r.queue].Push(r.readed)
		r.readed = ""
		return
	}

	q := r.queues[r.queue]
	last := q.Last()
	data := last.data + r.readed
	last.data = data[:r.queueEndIdx]
	q.End()

	r.readed = data[r.queueEndIdx:]
	r.queueStartIdx = -1
	r.queueEndIdx = -1
}

func (r *Router) startQueue() bool {
	var i int
	for name, q := range r.queues {
		i = q.StartIndex(r.readed)
		if i == -1 {
			continue
		}
		if i > r.queueStartIdx && r.queueStartIdx != -1 {
			continue
		}
		r.queueStartIdx = i
		r.queue = name
	}
	return r.queueStartIdx != -1
}

func (r *Router) endQueue() bool {
	queue := r.queues[r.queue]
	r.queueEndIdx = queue.EndIndex(queue.Last().data + r.readed)
	return r.queueEndIdx != -1
}
