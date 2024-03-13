package main

import (
	"os"

	"github.com/ejuju/my-ttyart-exhibit/internal/algolight"
	"github.com/ejuju/my-ttyart-exhibit/internal/gameoflife"
	"github.com/ejuju/my-ttyart-exhibit/internal/markode"
)

const (
	cmdGameOfLife = "game-of-life"
	cmdMarkode    = "markode"
	cmdAlgolight  = "algolight"
)

func main() {
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, cmdGameOfLife)
	}
	var run func() error
	switch os.Args[1] {
	default:
		panic("unknown command")
	case cmdGameOfLife:
		run = gameoflife.Run
	case cmdMarkode:
		run = markode.Run
	case cmdAlgolight:
		run = algolight.Run
	}
	err := run()
	if err != nil {
		panic(err)
	}
}
