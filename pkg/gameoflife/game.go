package gameoflife

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/ejuju/my-ttyart-exhibit/pkg/tty"
	"golang.org/x/term"
)

func Run() (err error) {
Base:
	for {
		// Hide terminal cursor and restore terminal state on exit.
		f := os.Stdout

		// Create new random grid.
		width, height, err := term.GetSize(int(f.Fd()))
		if err != nil {
			return fmt.Errorf("get terminal size: %w", err)
		}
		width = width / 2 // We use two-character of text width per cell.
		g := makeGrid(width, height)
		g.randomize(0)

		// Use terminal raw mode.
		state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("use terminal raw mode: %w", err)
		}
		defer term.Restore(int(os.Stdin.Fd()), state)

		// Read user input.
		input := make(chan byte, 1)
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

		// Run game loop.
		tui := tty.UI{Writer: f}
		tui.HideCursor()
		defer tui.ShowCursor()
		defer tui.ResetTextStyle()
		ticker := time.NewTicker(75 * time.Millisecond)
		for {
			select {
			case b := <-input:
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
				g.render(tui)
			}
		}
	}
}

type grid [][]bool

func (g grid) width() int  { return len(g) }
func (g grid) height() int { return len(g[0]) }

func makeGrid(width, height int) (g grid) {
	g = make([][]bool, width)
	for x := 0; x < g.width(); x++ {
		g[x] = make([]bool, height)
	}
	return g
}

func (g grid) randomize(seed int64) {
	randr := rand.New(rand.NewSource(seed))
	for x := 0; x < g.width(); x++ {
		for y := 0; y < g.height(); y++ {
			g[x][y] = randr.Int()%2 == 0
		}
	}
}

func (g grid) at(x, y int) bool {
	width, height := len(g), len(g[0])
	return g[(width+x)%width][(height+y)%height]
}

func (g grid) countNeighbours(x, y int) (count int) {
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

func (g *grid) update() (changed bool) {
	next := makeGrid(len(*g), len((*g)[0]))
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
	(*g) = next
	return changed
}

func (g grid) render(tui tty.UI) {
	for x := 0; x < g.width(); x++ {
		for y := 0; y < g.height(); y++ {
			tui.MoveTo(x*2, y)
			if alive := g[x][y]; alive {
				tui.SetBackgroundWhite()
			} else {
				tui.SetBackgroundBlack()
			}
			tui.Print("  ")
		}
	}
}
