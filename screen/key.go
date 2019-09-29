package screen

const (
	CharLineStart = 1
	CharBackward  = 2
	CharInterrupt = 3
	CharDelete    = 4
	CharLineEnd   = 5
	CharForward   = 6
	CharBell      = 7
	CharCtrlH     = 8
	CharTab       = 9
	CharCtrlJ     = 10
	CharKill      = 11
	CharCtrlL     = 12
	CharEnter     = 13
	CharNext      = 14
	CharPrev      = 16
	CharBckSearch = 18
	CharFwdSearch = 19
	CharTranspose = 20
	CharCtrlU     = 21
	CharCtrlW     = 23
	CharCtrlY     = 25
	CharCtrlZ     = 26
	CharEsc       = 27
	CharEscapeEx  = 91
	CharBackspace = 127

	esc    = "\033["
	Cursor = "\u2588"
)

var (
	ClearLine  = []byte(esc + "2K\r")
	MoveUp     = []byte(esc + "1A")
	MoveDown   = []byte(esc + "1B")
	HideCursor = esc + "?25l"
	ShowCursor = esc + "?25h"

	// KeyEnter is the default key for submission/selection.
	KeyEnter rune = CharEnter

	// KeyBackspace is the default key for deleting input text.
	KeyBackspace rune = CharBackspace

	// KeyPrev is the default key to go up during selection.
	KeyPrev        rune = CharPrev
	KeyPrevDisplay      = "↑"

	// KeyNext is the default key to go down during selection.
	KeyNext        rune = CharNext
	KeyNextDisplay      = "↓"

	// KeyBackward is the default key to page up during selection.
	KeyBackward        rune = CharBackward
	KeyBackwardDisplay      = "←"

	// KeyForward is the default key to page down during selection.
	KeyForward        rune = CharForward
	KeyForwardDisplay      = "→"

	// KeySearch is the default key to search
	KeySearch        rune = '/'
	KeySearchDisplay      = '/'
)
