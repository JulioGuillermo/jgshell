package shellinfrastructure

import "golang.org/x/sys/unix"

func configPty(fd uintptr) error {
	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return err
	}

	termios.Lflag |= unix.ECHO
	// termios.Lflag &^= unix.ECHO

	termios.Lflag |= unix.ICANON
	// termios.Lflag &^= unix.ICANON

	termios.Lflag |= unix.ISIG
	// termios.Lflag &^= unix.ISIG

	termios.Oflag |= unix.OPOST
	// termios.Oflag &^= unix.OPOST

	termios.Oflag |= unix.ONLCR
	// termios.Oflag &^= unix.ONLCR

	return unix.IoctlSetTermios(int(fd), unix.TCSETS, termios)
}
