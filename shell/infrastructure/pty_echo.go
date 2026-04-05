package shellinfrastructure

import "golang.org/x/sys/unix"

func configPty(fd uintptr) error {
	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return err
	}

	termios.Lflag |= unix.ECHO
	// termios.Lflag &^= unix.ECHO

	termios.Oflag |= unix.OPOST | unix.ONLCR
	termios.Lflag |= unix.ICANON | unix.ISIG | unix.ECHO

	return unix.IoctlSetTermios(int(fd), unix.TCSETS, termios)
}
