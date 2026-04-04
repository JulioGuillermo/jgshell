package executordomain

import "sync"

type Reader interface {
	Read(func(string) (string, bool)) (string, error)
	ReadPrecond(sync.Locker, func(string) bool, func(string) (string, bool)) (string, error)
}
