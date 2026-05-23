package pardon

import (
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/engmtcdrm/go-pardon/internal/keys"
)

type Input struct {
	// Hide indicates whether the input should be hidden (e.g., for password
	// input).
	Hide bool

	writer io.Writer
	reader *os.File
	result []rune
}

func NewInput() *Input {
	return &Input{
		Hide:   false,
		writer: os.Stdout,
		reader: os.Stdin,
	}
}

func (i *Input) rawReadline(f *os.File) (string, error) {
	var pending []byte
	i.result = nil // Reset results before reading new input

outer:
	for {
		var buf [8]byte
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 {
			if err == io.EOF {
				break outer
			}
			continue outer
		}

		pending = append(pending, buf[:n]...)

	inner:
		for len(pending) > 0 {
			if pending[0] == keys.CarriageReturn || pending[0] == keys.Enter {
				i.print("\n")
				return string(i.result), nil
			}

			if pending[0] == keys.CtrlC {
				return "", ErrUserAborted
			}

			if pending[0] == keys.Backspace || pending[0] == keys.Delete {
				pending = pending[1:]
				if len(i.result) > 0 {
					i.result = i.result[:len(i.result)-1]
					i.print("\b \b")
				}
				continue inner
			}

			if pending[0] == keys.Escape {
				if len(pending) < 2 {
					break inner
				}

				if pending[1] == keys.LeftBracket || pending[1] == keys.CapitalO {
					i := 2
				secondInner:
					for i < len(pending) {
						c := pending[i]
						if c >= 0x40 && c <= 0x7E {
							i++
							break secondInner
						}
						i++
					}

					if i > len(pending) {
						break inner
					}

					pending = pending[i:]
					continue inner
				}

				pending = pending[1:]
				continue inner
			}

			// pending, shouldContinue := i.handleEscapeSequence(pending)
			// if shouldContinue {
			// 	continue inner
			// } else {
			// 	break
			// }

			r, size := utf8.DecodeRune(pending)
			if r == utf8.RuneError && size == 1 {
				if !utf8.FullRune(pending) {
					break inner
				}

				pending = pending[1:]
				continue inner
			}

			pending = pending[size:]

			if unicode.IsPrint(r) && !unicode.IsControl(r) {
				i.result = append(i.result, r)
				i.print(string(r))
			}
		}
	}

	i.print("\n")
	return string(i.result), nil
}

func (i *Input) handleEscapeSequence(pending []byte) ([]byte, bool) {
	if len(pending) < 2 {
		return pending, false
	}

	if pending[0] != keys.Escape {
		return pending, false
	}

	// if pending[0] == keys.Escape {
	if pending[1] == keys.LeftBracket || pending[1] == keys.CapitalO {
		i := 2

		for i < len(pending) {
			c := pending[i]
			if c >= 0x40 && c <= 0x7E {
				i++
				break
			}
			i++
		}

		if i > len(pending) {
			return pending, false
		}

		pending = pending[i:]
		return pending, true
	}

	pending = pending[1:]
	return pending, true
	// }
}

// print writes the given arguments to the terminal if [Input.Hide] is false.
func (i *Input) print(a ...any) {
	if !i.Hide {
		fmt.Fprint(i.writer, a...)
	}
}
