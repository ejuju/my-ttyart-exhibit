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
	ui.EraseEntireScreen()
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

	generationTicker := time.NewTicker(time.Second / time.Duration(g.fps))
	timeoutTimer := time.NewTimer(5 * time.Minute)

	// Run game loop.
	for {
		select {
		case <-timeoutTimer.C:
			goto Base
		case b := <-input:
			switch b {
			default:
				goto Base
			case 'q':
				return nil
			case '+':
				if g.fps < 60 {
					g.fps++
				}
				generationTicker.Reset(time.Second / time.Duration(g.fps))
			case '-':
				if g.fps > 1 {
					g.fps--
				}
				generationTicker.Reset(time.Second / time.Duration(g.fps))
			}
		case <-generationTicker.C:
			g.update(ui)
		}
	}
}

type game struct {
	numRuns       int
	generation    int
	width, height int
	cells         []cell
	fps           int
}

type cell struct {
	x, y int
}

func randomCells(seed int64, width, height int) (cells []cell) {
	randr := rand.New(rand.NewSource(seed))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if randr.Int()%2 == 0 {
				cells = append(cells, cell{x: x, y: y})
			}
		}
	}
	return cells
}

func (g game) isAlive(x, y int) bool {
	x = (g.width + x) % g.width
	y = (g.height + y) % g.height
	for _, c := range g.cells {
		if c.x == x && c.y == y {
			return true
		}
	}
	return false
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

func (g *game) update(ui tty.TUI) {
	next := make([]cell, 0, g.width*g.height)
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			count := g.countNeighbours(x, y)
			isAliveNow := g.isAlive(x, y)
			isAliveNext := (isAliveNow && (count == 2 || count == 3)) || (!isAliveNow && count == 3)
			if isAliveNext {
				next = append(next, cell{x: x, y: y})
			}

			if isAliveNext {
				ui.SetBackgroundRGB(127, 0, 128+uint8((x+y+g.generation)%256)/2)
			} else {
				ui.SetBackgroundRGB(0, 0, 0)
			}
			ui.MoveTo(x*2, y)
			txt := "  "
			if count > 0 {
				txt = " " + strconv.Itoa(count)
			}
			ui.Print(txt)
		}
	}

	g.generation++
	g.cells = next

	// Render bottom banner.
	ui.SetBackgroundRGB(16, 16, 16)

	ui.MoveTo(0, g.height)
	content := fmt.Sprintf("#%d | %d FPS | Generation %d | Population %d", g.numRuns, g.fps, g.generation, len(g.cells))
	ui.Print(content + strings.Repeat(" ", g.width*2-len(content)))

	ui.MoveTo(0, g.height+1)
	content = "'+' = speed up | '-' = slow down | 'q' = quit | Press any other key to restart."
	ui.Print(content + strings.Repeat(" ", g.width*2-len(content)))

	ui.MoveTo(0, g.height+2)
	content = "Game of Life - John Conway"
	ui.Print(content + strings.Repeat(" ", g.width*2-len(content)))
}
