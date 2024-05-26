package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

type point struct {
	x, y int
}

const (
	width  = 30
	height = 20
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Failed to create the screen")

		os.Exit(1)
	}

	err = screen.Init()
	if err != nil {
		fmt.Printf("Failed to init the screen")

		os.Exit(1)
	}
	screen.Clear()

	game := Game{}
	game.Init(&screen)

	game.StartGame()
}
