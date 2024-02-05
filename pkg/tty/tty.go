package tty

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type TUI struct {
	beforeMakeRaw *term.State
	stdin, stdout *os.File
}

func NewTUI() (ui TUI) {
	return TUI{stdin: os.Stdin, stdout: os.Stdout}
}

func (ui TUI) Restore() (err error) {
	return term.Restore(int(ui.stdin.Fd()), ui.beforeMakeRaw)
}

func (ui TUI) Print(s string) {
	_, err := io.WriteString(ui.stdout, s)
	if err != nil {
		panic(err)
	}
}

func (ui TUI) Printf(f string, args ...any) {
	ui.Print(fmt.Sprintf(f, args...))
}

func (ui TUI) EraseEntireScreen()             { ui.Print("\x1b[2J") }
func (ui TUI) MoveTo(x, y int)                { ui.Printf("\x1b[%d;%dH", y+1, x+1) }
func (ui TUI) HideCursor()                    { ui.Print("\x1b[?25l") }
func (ui TUI) ShowCursor()                    { ui.Print("\x1b[?25h") }
func (ui TUI) ResetTextStyle()                { ui.Print("\x1b[0m") }
func (ui TUI) SetForegroundRGB(r, g, b uint8) { ui.Printf("\x1b[38;2;%d;%d;%dm", r, g, b) }
func (ui TUI) SetBackgroundRGB(r, g, b uint8) { ui.Printf("\x1b[48;2;%d;%d;%dm", r, g, b) }

func (ui TUI) Size() (width, height int, err error) {
	return term.GetSize(int(ui.stdout.Fd()))
}

func (ui *TUI) GoListenRaw(input chan byte) (err error) {
	ui.beforeMakeRaw, err = term.MakeRaw(int(ui.stdin.Fd()))
	if err != nil {
		return err
	}
	go func() {
		for {
			b := [1]byte{}
			_, err := io.ReadFull(os.Stdin, b[:])
			if err != nil {
				panic(err)
			}
			input <- b[0]
		}
	}()
	return nil
}
