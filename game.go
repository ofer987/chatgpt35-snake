package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type point struct {
	x, y int
}

type Snake struct {
	body      []point
	direction tcell.Key
	alive     bool
}

type Game struct {
	food       point
	score      int
	gameOver   bool
	Screen     *tcell.Screen
	event      Event
	snake      Snake
	gameStyle  tcell.Style
	snakeStyle tcell.Style
}

type Event interface {
}

type StateEvent struct {
	Exit bool
}

type MovementEvent struct {
	ChangeDirection tcell.Key
}

func (game *Game) Init(screen *tcell.Screen) {
	game.gameStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	game.snakeStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	game.Screen = screen

	game.snake = Snake{
		body:      []point{{5, 5}, {4, 5}, {3, 5}},
		direction: tcell.KeyRight,
		alive:     true,
	}

	game.placeFood()
}

func (game *Game) StartGame() {
	gEvent := make(chan Event)
	go game.inputLoop(gEvent)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case ev := <-gEvent:
			switch tev := ev.(type) {
			case StateEvent:
				if tev.Exit {
					return
				}
			case MovementEvent:
				game.changeDirection(tev.ChangeDirection)
			}
		case <-ticker.C:
			game.update()
			game.draw()
		}
	}
}

func (game *Game) inputLoop(gEvent chan<- Event) {
	for {
		ev := (*game.Screen).PollEvent()
		switch tev := ev.(type) {
		case *tcell.EventKey:
			handleInputEvent(tev, gEvent)
		}
	}
}

func handleInputEvent(ev *tcell.EventKey, gEvent chan<- Event) {
	switch ev.Key() {
	case tcell.KeyDown:
		fallthrough
	case 107:
		gEvent <- MovementEvent{tcell.KeyDown}
	case tcell.KeyUp:
		fallthrough
	case 106:
		gEvent <- MovementEvent{ChangeDirection: tcell.KeyUp}
	case tcell.KeyLeft:
		fallthrough
	case 108:
		gEvent <- MovementEvent{ChangeDirection: tcell.KeyLeft}
	case tcell.KeyRight:
		fallthrough
	case 104:
		gEvent <- MovementEvent{ChangeDirection: tcell.KeyRight}
	case tcell.KeyEnter:
		fallthrough
	case tcell.KeyEscape:
		gEvent <- StateEvent{true}
	}
}

func (game *Game) changeDirection(newDirection tcell.Key) {
	if !game.snake.alive {
		return
	}

	switch newDirection {
	case tcell.KeyLeft:
		if game.snake.direction == tcell.KeyUp || game.snake.direction == tcell.KeyDown {
			game.snake.direction = newDirection
		}
	case tcell.KeyRight:
		if game.snake.direction == tcell.KeyUp || game.snake.direction == tcell.KeyDown {
			game.snake.direction = newDirection
		}
	case tcell.KeyUp:
		if game.snake.direction == tcell.KeyLeft || game.snake.direction == tcell.KeyRight {
			game.snake.direction = newDirection
		}
	case tcell.KeyDown:
		if game.snake.direction == tcell.KeyLeft || game.snake.direction == tcell.KeyRight {
			game.snake.direction = newDirection
		}
	}
}

func (game *Game) update() {
	head := game.snake.body[0]
	var newHead point

	switch game.snake.direction {
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
		game.snake.alive = false
		return
	}

	for _, p := range game.snake.body[1:] {
		if newHead == p {
			game.snake.alive = false
			return
		}
	}

	// Check if the snake eats the food
	if newHead == game.food {
		game.score++
		game.placeFood()
	} else {
		game.snake.body = game.snake.body[:len(game.snake.body)-1]
	}

	game.snake.body = append([]point{newHead}, game.snake.body...)
}

func (game *Game) draw() {
	(*game.Screen).Clear()

	// Draw snake
	for i, p := range game.snake.body {
		char := ' '
		if i == 0 {
			char = '@' // Head
		} else {
			char = '#' // Body
		}

		(*game.Screen).SetContent(p.x, p.y, rune(char), nil, game.snakeStyle)
	}

	// Draw food
	(*game.Screen).SetContent(game.food.x, game.food.y, '$', nil, game.snakeStyle)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", game.score)
	for i, char := range scoreStr {
		(*game.Screen).SetContent(i, height+1, rune(char), nil, game.snakeStyle)
	}

	game.drawTopBorder()
	game.drawBottomBorder()
	game.drawLeftBorder()
	game.drawRightBorder()

	if !game.snake.alive {
		message := "You have lost"
		// (*game.Screen).SetContent(10, height+5, rune('k'), nil, game.snakeStyle)
		for i, char := range message {
			(*game.Screen).SetContent(i, height+2, rune(char), nil, game.snakeStyle)
		}
	}

	(*game.Screen).Show()
}

func (game *Game) drawLeftBorder() {
	for y := 1; y < height-1; y += 1 {
		(*game.Screen).SetContent(0, y, '|', nil, game.snakeStyle)
	}
}

func (game *Game) drawRightBorder() {
	for y := 1; y < height-1; y += 1 {
		(*game.Screen).SetContent(width, y, '|', nil, game.snakeStyle)
	}
}

func (game *Game) drawTopBorder() {
	// Top-Left corner
	(*game.Screen).SetContent(0, 0, '/', nil, game.snakeStyle)

	for x := 1; x < width; x += 1 {
		(*game.Screen).SetContent(x, 0, '-', nil, game.snakeStyle)
	}

	// Top-Right corner
	(*game.Screen).SetContent(width, 0, '\\', nil, game.snakeStyle)
}

func (game *Game) drawBottomBorder() {
	// Bottom-Left corner
	(*game.Screen).SetContent(0, height-1, '\\', nil, game.snakeStyle)

	for x := 1; x < width; x += 1 {
		(*game.Screen).SetContent(x, height-1, '-', nil, game.snakeStyle)
	}

	// Bottom-Right corner
	(*game.Screen).SetContent(width, height-1, '/', nil, game.snakeStyle)
}

func (game *Game) placeFood() {
	game.food = point{1 + rand.Intn(width-2), 1 + rand.Intn(height-2)}
}