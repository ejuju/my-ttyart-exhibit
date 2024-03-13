// Markov-quine.
package markode

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
	"unicode/utf8"

	_ "embed"

	"github.com/ejuju/my-ttyart-exhibit/pkg/tty"
)

const order = 3

//go:embed matrix.txt
var srcCode string

func Run() (err error) {
	// Hide terminal cursor and restore terminal state on exit.
	ui := tty.NewTUI()
	defer ui.ShowCursor()
	defer ui.ResetTextStyle()
	ui.HideCursor()
	ui.ResetTextStyle()
	ui.MoveTo(0, 0)
	ui.EraseEntireScreen()

	mc := newMarkovChain(srcCode)
	start := [order]rune([]rune(srcCode[:order]))
	g := game{
		charsPerSec: 100,
		ran:         rand.New(rand.NewSource(time.Now().UnixNano())),
		mc:          mc,
		lastchars:   start,
	}

	inputs := make(chan string)
	srv := NewServer(inputs)
	go func() {
		err := http.ListenAndServe(":8080", srv)
		if err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ticker := time.NewTicker(time.Second / time.Duration(g.charsPerSec))
	for {
		select {
		case txt := <-inputs:
			mc.learn([]rune(txt))
		case <-interrupt:
			return nil
		case <-ticker.C:
			g.tick(ui)
		}
	}
}

type game struct {
	ran         *rand.Rand
	mc          markovChain
	lastchars   [order]rune
	charsPerSec int
}

func (g *game) tick(ui tty.TUI) {
	c := g.mc.next(g.ran, string(g.lastchars[:]))
	ui.Print(string(c))
	for i := 1; i < order; i++ {
		g.lastchars[i-1] = g.lastchars[i]
	}
	g.lastchars[len(g.lastchars)-1] = c
}

// maps n-grams to a "next" character.
type markovChain map[string][]rune

func newMarkovChain(corpus string) (mc markovChain) {
	mc = markovChain{}
	mc.learn([]rune(corpus))
	return mc
}

func (mc *markovChain) learn(corpus []rune) {
	for i := 0; i < len(corpus)-order; i++ {
		currSequence := string(corpus[i : i+order])
		(*mc)[currSequence] = append((*mc)[currSequence], corpus[i+order])
	}
}

func (mc markovChain) next(ran *rand.Rand, str string) rune {
	strNumChars := utf8.RuneCountInString(str)
	if strNumChars < order {
		panic(fmt.Errorf("input string %d is shorter than order %d", strNumChars, order))
	}
	lastNchars := []rune(str)[strNumChars-order:]
	options, ok := mc[string(lastNchars)]
	if !ok {
		for _, v := range mc {
			return v[rand.Intn(len(v))]
		}
	}
	return options[ran.Intn(len(options))]
}
