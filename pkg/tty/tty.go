package tty

import (
	"fmt"
	"io"
)

type UI struct{ io.Writer }

func (ui UI) Print(s string) {
	_, err := io.WriteString(ui, s)
	if err != nil {
		panic(err)
	}
}

func (ui UI) Printf(f string, args ...any) {
	ui.Print(fmt.Sprintf(f, args...))
}

func (ui UI) MoveTo(x, y int)     { ui.Printf("\x1b[%d;%dH", y+1, x+1) }
func (ui UI) HideCursor()         { ui.Print("\x1b[?25l") }
func (ui UI) ShowCursor()         { ui.Print("\x1b[?25h") }
func (ui UI) ResetTextStyle()     { ui.Print("\x1b[0m") }
func (ui UI) SetBackgroundBlack() { ui.Print("\x1b[40m") }
func (ui UI) SetBackgroundWhite() { ui.Print("\x1b[47m") }
