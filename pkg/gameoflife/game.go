package gameoflife

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ejuju/my-ttyart-exhibit/pkg/tty"
)

func Run() (err error) {
	// Hide terminal cursor and restore terminal state on exit.
	ui := tty.NewTUI()
	ui.HideCursor()
	defer ui.ShowCursor()
	defer ui.ResetTextStyle()

	// Use terminal raw mode.
	input := make(chan byte)
	err = ui.GoListenRaw(input)
	if err != nil {
		return fmt.Errorf("use terminal raw mode: %w", err)
	}
	defer ui.Restore()

	g := game{fps: 10}

Base:
	// Create new random grid depending on terminal size.
	g.width, g.height, err = ui.Size()
	if err != nil {
		return fmt.Errorf("get terminal size: %w", err)
	}
	g.width = g.width / 2   // We use two-character of text width per cell.
	g.height = g.height - 3 // We use the bottom lines for banner.

	ui.ResetTextStyle()
	ui.MoveTo(0, 0)
	ui.EraseEntireScreen()

	g.numRuns++
	g.generation = 0
	g.cells = randomCells(time.Now().UnixNano(), g.width, g.height)

	ticker := time.NewTicker(time.Second / time.Duration(g.fps))
	timeout := time.NewTimer(5 * time.Minute)

	// Run game loop.
	for {
		select {
		case <-timeout.C:
			goto Base
		case b := <-input:
			switch {
			default:
				goto Base
			case b == 'q':
				return nil
			case b == '+' && g.fps < 60:
				g.fps++
				ticker.Reset(time.Second / time.Duration(g.fps))
			case b == '-' && g.fps > 1:
				g.fps--
				ticker.Reset(time.Second / time.Duration(g.fps))
			}
		case <-ticker.C:
			g.tick(ui)
		}
	}
}

type game struct {
	numRuns       int
	generation    int
	width, height int
	cells         []bool
	fps           int
}

func randomCells(seed int64, width, height int) (cells []bool) {
	randr := rand.New(rand.NewSource(seed))
	cells = make([]bool, width*height)
	for i := range cells {
		cells[i] = randr.Int()%3 == 0
	}
	return cells
}

func (g game) isAlive(x, y int) bool {
	x = (g.width + x) % g.width
	y = (g.height + y) % g.height
	return g.cells[g.width*y+x]
}

func (g game) countNeighbours(x, y int) (count int) {
	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			if i == y && j == x {
				continue // Ignore own position.
			} else if g.isAlive(j, i) {
				count++
			}
		}
	}
	return count
}

func (g *game) tick(ui tty.TUI) {
	next := make([]bool, g.width*g.height)
	population := 0
	for i, isAliveNow := range g.cells {
		x, y := (i % g.width), (i / g.width)
		count := g.countNeighbours(x, y)
		isAliveNext := (isAliveNow && (count == 2 || count == 3)) || (!isAliveNow && count == 3)
		if isAliveNext {
			population++
		}
		next[g.width*y+x] = isAliveNext

		txt := "  "
		if isAliveNext {
			ui.SetBackgroundRGB(0, 0, 0)
		} else {
			txt = " " + strconv.Itoa(count)
			ui.SetBackgroundRGB(0xBB, 0x11, 0xFF)
		}
		ui.MoveTo(x*2, y)
		ui.Print(txt)
	}

	g.generation++
	g.cells = next

	// Render bottom banner.
	ui.SetBackgroundRGB(16, 16, 16)

	ui.MoveTo(0, g.height)
	content := fmt.Sprintf("#%d | %d FPS | Generation %d | Population %d", g.numRuns, g.fps, g.generation, population)
	ui.Print(content + strings.Repeat(" ", g.width*2-len(content)))

	ui.MoveTo(0, g.height+1)
	content = "'+' = speed up | '-' = slow down | 'q' = quit | Press any other key to restart."
	ui.Print(content + strings.Repeat(" ", g.width*2-len(content)))

	ui.MoveTo(0, g.height+2)
	content = "Game of Life - John Conway"
	ui.Print(content + strings.Repeat(" ", g.width*2-len(content)))
}
