package shellinfrastructure

import (
	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/unix"
)

func configPty(fd uintptr) error {
	_, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}

	termios, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	if err != nil {
		return err
	}
	// termios.Iflag |= unix.ICRNL

	termios.Lflag |= unix.ECHO
	// termios.Lflag &^= unix.ECHO

	// termios.Lflag |= unix.ICANON
	// termios.Lflag &^= unix.ICANON

	termios.Lflag |= unix.ISIG
	// termios.Lflag &^= unix.ISIG

	// termios.Oflag |= unix.OPOST
	// termios.Oflag &^= unix.OPOST

	// termios.Oflag |= unix.ONLCR
	// termios.Oflag &^= unix.ONLCR

	return unix.IoctlSetTermios(int(fd), unix.TCSETS, termios)
}
