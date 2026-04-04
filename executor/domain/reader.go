package executordomain

type Reader interface {
	Read(func(string) (string, bool)) (string, error)
}
