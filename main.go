package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	width  = 30
	height = 20
)

type point struct {
	x, y int
}

type Snake struct {
	body      []point
	direction tcell.Key
	alive     bool
}

var (
	food       point
	score      int
	gameOver   bool
	gameScreen tcell.Screen
	snakeStyle tcell.Style
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

	gameScreen = screen
	gameScreen.HideCursor()
	gameScreen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	snakeStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	gameScreen.Clear()

	snake := Snake{
		body:      []point{{5, 5}, {4, 5}, {3, 5}},
		direction: tcell.KeyRight,
		alive:     true,
	}

	placeFood()

	ge := make(chan gameEvent)

	go gameLoop(ge)

	displayLoop(&snake, ge)
}

func gameLoop(gEvent chan<- gameEvent) {
	for {
		ev := gameScreen.PollEvent()
		switch tev := ev.(type) {
		case *tcell.EventKey:
			handleGameEvent(tev, gEvent)
		}
	}
}

type gameEvent struct {
	Exit            bool
	ChangeDirection tcell.Key
}

func handleGameEvent(ev *tcell.EventKey, gEvent chan<- gameEvent) {
	switch ev.Key() {
	case tcell.KeyDown:
		fallthrough
	case 107:
		gEvent <- gameEvent{Exit: false, ChangeDirection: tcell.KeyDown}
	case tcell.KeyUp:
		fallthrough
	case 106:
		gEvent <- gameEvent{Exit: false, ChangeDirection: tcell.KeyUp}
	case tcell.KeyLeft:
		fallthrough
	case 108:
		gEvent <- gameEvent{Exit: false, ChangeDirection: tcell.KeyLeft}
	case tcell.KeyRight:
		fallthrough
	case 104:
		gEvent <- gameEvent{Exit: false, ChangeDirection: tcell.KeyRight}
	case tcell.KeyEnter:
		fallthrough
	case tcell.KeyEscape:
		gEvent <- gameEvent{true, tcell.KeyDown}
	}
}

func displayLoop(snake *Snake, gEvent <-chan gameEvent) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		draw(snake)

		select {
		case ev := <-gEvent:
			if ev.Exit {
				return
			}

			changeDirection(snake, ev.ChangeDirection)
		case <-ticker.C:
			update(snake)
		}
	}
}

func changeDirection(snake *Snake, newDirection tcell.Key) {
	switch newDirection {
	case tcell.KeyLeft:
		if snake.direction == tcell.KeyUp || snake.direction == tcell.KeyDown {
			snake.direction = newDirection
		}
	case tcell.KeyRight:
		if snake.direction == tcell.KeyUp || snake.direction == tcell.KeyDown {
			snake.direction = newDirection
		}
	case tcell.KeyUp:
		if snake.direction == tcell.KeyLeft || snake.direction == tcell.KeyRight {
			snake.direction = newDirection
		}
	case tcell.KeyDown:
		if snake.direction == tcell.KeyLeft || snake.direction == tcell.KeyRight {
			snake.direction = newDirection
		}
	}
}

func update(s *Snake) {
	head := s.body[0]
	var newHead point

	switch s.direction {
	case tcell.KeyUp:
		fallthrough
	case 107:
		newHead = point{head.x, head.y - 1}
	case tcell.KeyDown:
		fallthrough
	case 106:
		newHead = point{head.x, head.y + 1}
	case tcell.KeyLeft:
		fallthrough
	case 108:
		newHead = point{head.x - 1, head.y}
	case tcell.KeyRight:
		fallthrough
	case 104:
		newHead = point{head.x + 1, head.y}
	}

	// Check if the snake collides with the walls or itself
	if newHead.x <= 0 || newHead.x >= width || newHead.y <= 0 || newHead.y >= height-1 {
		s.alive = false
		return
	}

	for _, p := range s.body[1:] {
		if newHead == p {
			s.alive = false
			return
		}
	}

	// Check if the snake eats the food
	if newHead == food {
		score++
		placeFood()
	} else {
		s.body = s.body[:len(s.body)-1]
	}

	s.body = append([]point{newHead}, s.body...)
}

func draw(s *Snake) {
	gameScreen.Clear()

	// Draw snake
	for i, p := range s.body {
		char := ' '
		if i == 0 {
			char = '@' // Head
		} else {
			char = '#' // Body
		}

		gameScreen.SetContent(p.x, p.y, rune(char), nil, snakeStyle)
	}

	// Draw food
	gameScreen.SetContent(food.x, food.y, '$', nil, snakeStyle)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", score)
	for i, char := range scoreStr {
		gameScreen.SetContent(i, height+5, rune(char), nil, snakeStyle)
	}

	drawTopBorder()
	drawBottomBorder()
	drawLeftBorder()
	drawRightBorder()

	if !s.alive {
		message := "You have lost"

		for i, char := range message {
			gameScreen.SetContent(i, height+5, rune(char), nil, snakeStyle)
		}
	}

	gameScreen.Show()
}

func drawLeftBorder() {
	for y := 1; y < height-1; y += 1 {
		gameScreen.SetContent(0, y, '|', nil, snakeStyle)
	}
}

func drawRightBorder() {
	for y := 1; y < height-1; y += 1 {
		gameScreen.SetContent(width, y, '|', nil, snakeStyle)
	}
}

func drawTopBorder() {
	// Top-Left corner
	gameScreen.SetContent(0, 0, '/', nil, snakeStyle)

	for x := 1; x < width; x += 1 {
		gameScreen.SetContent(x, 0, '-', nil, snakeStyle)
	}

	// Top-Right corner
	gameScreen.SetContent(width, 0, '\\', nil, snakeStyle)
}

func drawBottomBorder() {
	// Bottom-Left corner
	gameScreen.SetContent(0, height-1, '\\', nil, snakeStyle)

	for x := 1; x < width; x += 1 {
		gameScreen.SetContent(x, height-1, '-', nil, snakeStyle)
	}

	// Bottom-Right corner
	gameScreen.SetContent(width, height-1, '/', nil, snakeStyle)
}

func placeFood() {
	food = point{1 + rand.Intn(width-2), 1 + rand.Intn(height-2)}
}
