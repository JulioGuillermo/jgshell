package executorapplication

import "strings"

func (r *Reader) Clear(str string) string {
	if str == "" {
		return str
	}

	src := []byte(str)
	var out strings.Builder
	out.Grow(len(src))

	for i := 0; i < len(src); {
		b := src[i]

		if b == 0x1b { // ESC
			if i+1 >= len(src) {
				i++
				continue
			}

			switch src[i+1] {
			case '[': // CSI: ESC [ ... final-byte(0x40-0x7E)
				i += 2
				for i < len(src) {
					c := src[i]
					i++
					if c >= 0x40 && c <= 0x7E {
						break
					}
				}
				continue

			case ']': // OSC: ESC ] ... (BEL or ESC \)
				start := i
				i += 2
				payloadStart := i
				termLen := 0

				for i < len(src) {
					if src[i] == 0x07 { // BEL
						termLen = 1
						break
					}
					if src[i] == 0x1b && i+1 < len(src) && src[i+1] == '\\' { // ST
						termLen = 2
						break
					}
					i++
				}

				if termLen == 0 {
					// Incomplete OSC sequence in current buffer chunk: drop it.
					break
				}

				payload := string(src[payloadStart:i])
				end := i + termLen

				if strings.HasPrefix(payload, "JGSHELL;") {
					out.Write(src[start:end])
				}

				i = end
				continue

			default:
				// Other ESC sequences (single/two-byte control forms): drop.
				i += 2
				continue
			}
		}

		// Remove ASCII control chars except LF/CR/TAB.
		if b < 0x20 {
			if b == '\n' || b == '\r' || b == '\t' {
				out.WriteByte(b)
			}
			i++
			continue
		}
		if b == 0x7f { // DEL
			i++
			continue
		}

		out.WriteByte(b)
		i++
	}

	return out.String()
}
