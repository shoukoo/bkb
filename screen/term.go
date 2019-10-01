package screen

import (
	"bytes"
	"fmt"
	"io"
)

type Terminal struct {
	w      io.Writer
	buf    *bytes.Buffer
	reset  bool
	cursor int
	height int
}

func New(w io.Writer) *Terminal {
	return &Terminal{buf: &bytes.Buffer{}, w: w}
}

// WriteString is a convenient function to write a new line passing a string.
// Check ScreenBuf.Write() for a detailed explanation of the function behaviour.
func (t *Terminal) WriteString(str string) (int, error) {
	return t.Write([]byte(str))
}

// Write writes a single line to the underlying buffer. If the Terminal was
// previously reset, all previous lines are clears and the output starts from
// the top. Line with \r and \n will cause an error since they can intefere
// with the termianl ability to Move between line
func (t *Terminal) Write(b []byte) (int, error) {

	for bytes.ContainsAny(b, "\r\n") {
		return 0, fmt.Errorf("%q should not contain either \\r or \\n", b)
	}

	if t.reset {
		for i := 0; i < t.height; i++ {
			_, err := t.buf.Write(MoveUp)
			if err != nil {
				return 0, err
			}
			_, err = t.buf.Write(ClearLine)
			if err != nil {
				return 0, err
			}
		}
		t.cursor = 0
		t.height = 0
		t.reset = false
	}

	switch {
	case t.cursor == t.height:
		n, err := t.buf.Write(ClearLine)
		if err != nil {
			return n, err
		}
		n, err = t.buf.Write(b)
		if err != nil {
			return n, err
		}
		n, err = t.buf.Write([]byte("\n"))
		if err != nil {
			return n, err
		}
		t.height++
		t.cursor++
		return n, nil
	case t.cursor < t.height:
		n, err := t.buf.Write(ClearLine)
		if err != nil {
			return n, err
		}
		n, err = t.buf.Write(b)
		if err != nil {
			return n, err
		}
		n, err = t.buf.Write(MoveDown)
		if err != nil {
			return n, err
		}
		t.cursor++
		return n, nil
	default:
		return 0, fmt.Errorf("invalid write cursor position (%d) exceeded line height: %d", t.cursor, t.height)
	}
}

// Flush writes any buffered data to the underlying io.Writer, ensuring that any pending data is displayed.
func (t *Terminal) Flush() error {
	for i := t.cursor; i < t.height; i++ {
		if i < t.height {
			_, err := t.buf.Write(ClearLine)
			if err != nil {
				return err
			}
		}
		_, err := t.buf.Write(MoveDown)
		if err != nil {
			return err
		}
	}

	_, err := t.buf.WriteTo(t.w)
	if err != nil {
		return err
	}

	t.buf.Reset()

	for i := 0; i < t.height; i++ {
		_, err := t.buf.Write(MoveUp)
		if err != nil {
			return err
		}
	}

	t.cursor = 0

	return nil
}
