package main

import (
	"os"

	"github.com/ejuju/my-ttyart-exhibit/pkg/gameoflife"
	"github.com/ejuju/my-ttyart-exhibit/pkg/markode"
)

const (
	cmdGameOfLife = "game-of-life"
	cmdMarkode    = "markode"
)

func main() {
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, cmdGameOfLife)
	}
	var run func() error
	switch os.Args[1] {
	case cmdGameOfLife:
		run = gameoflife.Run
	case cmdMarkode:
		run = markode.Run
	}
	err := run()
	if err != nil {
		panic(err)
	}
}
