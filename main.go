package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	width  = 30
	height = 20
)

func main() {
	var screen tcell.Screen
	var err error

	if screen, err = tcell.NewScreen(); err != nil {
		fmt.Printf("Failed to ,create the screen")

		os.Exit(1)
	}

	if err = screen.Init(); err != nil {
		fmt.Printf("Failed to init the screen")

		os.Exit(1)
	}
	screen.Clear()

	game := Game{}
	game.Init(&screen)

	game.StartGame()
}
