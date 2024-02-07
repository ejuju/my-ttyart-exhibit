package algolight

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/ejuju/my-ttyart-exhibit/pkg/tty"
)

func Run() (err error) {
	ui := tty.NewTUI()
	defer ui.ShowCursor()
	defer ui.ResetTextStyle()
	ui.HideCursor()
	ui.ResetTextStyle()
	ui.MoveTo(0, 0)
	ui.EraseEntireScreen()

	g := grid{}
	g.width, g.height, err = ui.Size()
	if err != nil {
		return fmt.Errorf("get terminal size: %w", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ticker := time.NewTicker(time.Second / time.Duration(24)) // 24 FPS
	for {
		select {
		case <-interrupt:
			return nil
		case <-ticker.C:
			g.tick(ui)
		}
	}
}

type grid struct {
	width, height int
	ticks         int
}

func (g *grid) tick(ui tty.TUI) {
	cyclic := uint8(0)

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			if y == 0 {
				continue
			}
			ui.MoveTo(x, y)
			amplitude := 155
			cyclic = uint8(math.Abs(float64((g.ticks % (amplitude * 2)) - amplitude)))
			ui.SetBackgroundRGB(uint8(x), uint8(y), cyclic)
			ui.Print(" ")
		}
	}

	g.ticks++
}
