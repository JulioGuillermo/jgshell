package executorapplication

import "strings"

const (
	esc = byte(0x1b)
	bel = byte(0x07)
	del = byte(0x7f)
)

type Cleanner struct {
	str string
	src []byte
	out strings.Builder
}

func NewCleanner() *Cleanner {
	return &Cleanner{}
}

func (c *Cleanner) Clear(str string) string {
	if str == "" {
		return ""
	}

	c.str = str
	c.src = []byte(c.str)

	if cut := c.findTrailingIncompleteEscape(); cut >= 0 {
		c.str = c.str[:cut]
		c.src = []byte(c.str)
	}
	if len(c.src) == 0 {
		return ""
	}

	c.out.Reset()
	c.out.Grow(len(c.src))

	for i := 0; i < len(c.src); {
		b := c.src[i]

		if b == esc {
			i = c.consumeEscapeSequence(i)
			continue
		}

		if c.shouldDropControl(b) {
			i++
			continue
		}

		c.out.WriteByte(b)
		i++
	}

	return c.out.String()
}

func (c *Cleanner) findTrailingIncompleteEscape() int {
	for i := 0; i < len(c.src); {
		if c.src[i] != esc {
			i++
			continue
		}

		start := i
		next, incomplete := c.consumeEscapeForCut(i)
		if incomplete {
			return start
		}
		i = next
	}
	return -1
}

func (c *Cleanner) consumeEscapeForCut(i int) (int, bool) {
	if i+1 >= len(c.src) {
		return i, true
	}

	switch c.src[i+1] {
	case '[':
		return c.consumeCSIForCut(i)
	case ']':
		return c.consumeOSCForCut(i)
	default:
		return i + 2, false
	}
}

func (c *Cleanner) consumeCSIForCut(i int) (int, bool) {
	i += 2
	for i < len(c.src) {
		ch := c.src[i]
		i++
		if ch >= 0x40 && ch <= 0x7E {
			return i, false
		}
	}
	return 0, true
}

func (c *Cleanner) consumeOSCForCut(i int) (int, bool) {
	i += 2
	for i < len(c.src) {
		if c.src[i] == bel {
			return i + 1, false
		}
		if c.src[i] == esc {
			if i+1 >= len(c.src) {
				return 0, true
			}
			if c.src[i+1] == '\\' {
				return i + 2, false
			}
		}
		i++
	}
	return 0, true
}

func (c *Cleanner) consumeEscapeSequence(i int) int {
	if i+1 >= len(c.src) {
		return i + 1
	}

	switch c.src[i+1] {
	case '[':
		return c.consumeCSI(i)
	case ']':
		return c.consumeOSC(i)
	default:
		// Other ESC sequences are ignored.
		return i + 2
	}
}

func (c *Cleanner) consumeCSI(i int) int {
	// CSI format: ESC [ ... final-byte(0x40-0x7E)
	i += 2
	for i < len(c.src) {
		ch := c.src[i]
		i++
		if ch >= 0x40 && ch <= 0x7E {
			break
		}
	}
	return i
}

func (c *Cleanner) consumeOSC(i int) int {
	// OSC format: ESC ] payload (BEL | ESC \)
	start := i
	payloadStart := i + 2
	termStart, termLen := c.findOSCTerminator(payloadStart)
	if termLen == 0 {
		// Incomplete OSC sequence in current chunk: drop remaining tail.
		return len(c.src)
	}

	payload := string(c.src[payloadStart:termStart])
	end := termStart + termLen
	if strings.HasPrefix(payload, "JGSHELL;") {
		c.out.Write(c.src[start:end])
	}

	return end
}

func (c *Cleanner) findOSCTerminator(from int) (int, int) {
	for i := from; i < len(c.src); i++ {
		if c.src[i] == bel {
			return i, 1
		}
		if c.src[i] == esc && i+1 < len(c.src) && c.src[i+1] == '\\' {
			return i, 2
		}
	}
	return 0, 0
}

func (c *Cleanner) shouldDropControl(b byte) bool {
	if b == '\n' || b == '\r' || b == '\t' {
		return false
	}
	if b == del {
		return true
	}
	return b < 0x20
}
