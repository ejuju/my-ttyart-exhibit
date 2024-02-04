package gameoflife

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ejuju/my-ttyart-exhibit/pkg/tty"
)

func Run() (err error) {
	// Hide terminal cursor and restore terminal state on exit.
	tui := tty.NewTUI()
	tui.HideCursor()
	tui.EraseEntireScreen()
	defer tui.ShowCursor()
	defer tui.ResetTextStyle()

	// Use terminal raw mode.
	input := make(chan byte)
	err = tui.GoListenRaw(input)
	if err != nil {
		return fmt.Errorf("use terminal raw mode: %w", err)
	}
	defer tui.Restore()

	// Create new random grid.
	width, height, err := tui.Size()
	if err != nil {
		return fmt.Errorf("get terminal size: %w", err)
	}
	width = width / 2   // We use two-character of text width per cell.
	height = height - 2 // We use the bottom two lines for banner.
	g := game{cells: makeCells(width, height-2)}

Base:
	// Run game loop.
	for {
		g.numRuns++
		g.generation = 0
		g.randomize(time.Now().UnixNano())
		g.render(tui)

		ticker := time.NewTicker(100 * time.Millisecond)
		previousPopulationCounts := [2]int{}
		for {
			select {
			case b := <-input:
				tui.ResetTextStyle()
				tui.MoveTo(0, 0)
				tui.EraseEntireScreen()
				switch b {
				default:
					continue Base
				case 'q':
					return nil
				}
			case <-ticker.C:
				changed := g.update()
				if !changed {
					return nil
				}
				population := g.population()
				isStuck := previousPopulationCounts[0] == population && previousPopulationCounts[1] == population
				if isStuck {
					continue Base
				}
				previousPopulationCounts[0] = previousPopulationCounts[1]
				previousPopulationCounts[1] = population

				g.render(tui)
			}
		}
	}
}

type game struct {
	numRuns    int
	generation int
	cells      [][]bool
}

func (g game) width() int  { return len(g.cells) }
func (g game) height() int { return len(g.cells[0]) }

func makeCells(width, height int) (cells [][]bool) {
	cells = make([][]bool, width)
	for x := 0; x < width; x++ {
		cells[x] = make([]bool, height)
	}
	return cells
}

func (g game) randomize(seed int64) {
	randr := rand.New(rand.NewSource(seed))
	for x := 0; x < g.width(); x++ {
		for y := 0; y < g.height(); y++ {
			g.cells[x][y] = randr.Int()%2 == 0
		}
	}
}

func (g game) at(x, y int) bool {
	width, height := g.width(), g.height()
	return g.cells[(width+x)%width][(height+y)%height]
}

func (g game) countNeighbours(x, y int) (count int) {
	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			if i == y && j == x {
				continue // Ignore own position.
			} else if g.at(j, i) {
				count++
			}
		}
	}
	return count
}

func (g game) population() (count int) {
	for x := 0; x < g.width(); x++ {
		for y := 0; y < g.height(); y++ {
			if g.at(x, y) {
				count++
			}
		}
	}
	return count
}

func (g *game) update() (changed bool) {
	next := makeCells(g.width(), g.height())
	for x := 0; x < g.width(); x++ {
		for y := 0; y < g.height(); y++ {
			count := g.countNeighbours(x, y)
			alive := g.at(x, y)
			if (alive && (count == 2 || count == 3)) || (!alive && count == 3) {
				next[x][y] = true
				changed = true
			}
		}
	}
	g.generation++
	g.cells = next
	return changed
}

func (g game) render(tui tty.TUI) {
	for x := 0; x < g.width(); x++ {
		for y := 0; y < g.height(); y++ {
			tui.MoveTo(x*2, y)
			if alive := g.at(x, y); alive {
				tui.SetBackgroundRGB(127, 0, 128+uint8((x+y+g.generation)%256)/2)
			} else {
				tui.SetBackgroundRGB(0, 0, 0)
			}
			tui.Print("  ")
		}
	}

	tui.SetBackgroundRGB(0, 0, 0)

	tui.MoveTo(0, g.height()+0)
	tui.Printf("[#%d] Generation %d | Population %d", g.numRuns, g.generation, g.population())

	tui.MoveTo(0, g.height()+1)
	tui.Print("Press 'q' to quit, press any other key to restart.")

	tui.MoveTo(0, g.height()+2)
	tui.Print("Game of Life - John Conway.")
}
