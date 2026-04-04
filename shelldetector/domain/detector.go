package shelldetectordomain

type ShellDetector interface {
	DetectShell() (string, error)
}
