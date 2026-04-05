package shellinfrastructure

import "golang.org/x/sys/unix"

func setEcho(fd uintptr, on bool) error {
	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return err
	}

	if on {
		termios.Lflag |= unix.ECHO
	} else {
		termios.Lflag &^= unix.ECHO
	}

	return unix.IoctlSetTermios(int(fd), unix.TCSETS, termios)
}
