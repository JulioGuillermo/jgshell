package routerapplication

import (
	"regexp"

	routerdomain "github.com/julioguillermo/jgshell/router/domain"
)

const MaxQueueSize = 30

type Queue struct {
	buffer chan *Element
	last   *Element
	name   string
	start  *regexp.Regexp
	end    *regexp.Regexp
}

func NewQueue(name, start, end string) *Queue {
	return &Queue{
		buffer: make(chan *Element, MaxQueueSize),
		name:   name,
		start:  regexp.MustCompile(start),
		end:    regexp.MustCompile(end),
	}
}

func (q *Queue) Name() string {
	return q.name
}

func (q *Queue) Last() *Element {
	return q.last
}

func (q *Queue) Start() {
	q.last = NewElement()
	q.buffer <- q.last
}

func (q *Queue) End() {
	q.Last().Close()
}

func (q *Queue) Push(data string) {
	last := q.Last()
	last.AppendData(data)
}

func (q *Queue) Pop() routerdomain.Element {
	return <-q.buffer
}

func (q *Queue) Clear() {
	q.buffer = make(chan *Element, MaxQueueSize)
}

func (q *Queue) StartIndex(data string) int {
	loc := q.start.FindStringIndex(data)
	if loc == nil {
		return -1
	}
	return loc[0]
}

func (q *Queue) EndIndex(data string) int {
	loc := q.end.FindStringIndex(data)
	if loc == nil {
		return -1
	}
	return loc[1]
}
