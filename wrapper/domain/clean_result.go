package wrapperdomain

import "time"

type CmdUnwrapResult struct {
	Output string

	Started   bool
	IsRunning bool

	User    string
	Pwd     string
	Code    int
	EndTime *time.Time
}
