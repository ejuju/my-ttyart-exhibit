package main

import "github.com/ejuju/my-ttyart-exhibit/pkg/gameoflife"

func main() {
	err := gameoflife.Run()
	if err != nil {
		panic(err)
	}
}
