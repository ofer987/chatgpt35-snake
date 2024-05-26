package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameStatus uint16

const (
	Started GameStatus = iota
	Success
	Failed
)

type Game struct {
	food        point
	score       int
	status      GameStatus
	Screen      *tcell.Screen
	event       Event
	snake       Snake
	borderStyle tcell.Style
	snakeStyle  tcell.Style
	headStyle   tcell.Style
	bodyStyle   tcell.Style
	foodStyle   tcell.Style
	textStyle   tcell.Style
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
	game.borderStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	game.snakeStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	game.headStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlueViolet)
	game.bodyStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
	game.foodStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	game.textStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	game.Screen = screen

	game.snake = CreateSnake()

	game.placeFood()
}

func (game *Game) StartGame() {
	game.status = Started

	gEvent := make(chan Event)
	go game.inputLoop(gEvent)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer (*game.Screen).Fini()

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
			game.handleInputEvent(tev, gEvent)
		}
	}
}

func (game *Game) handleInputEvent(ev *tcell.EventKey, gEvent chan<- Event) {
	if game.snake.movement == Changed {
		return
	}

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
	if game.status != Started {
		return
	}

	switch newDirection {
	case tcell.KeyUp:
		game.snake.MoveUp()
	case tcell.KeyRight:
		game.snake.MoveRight()
	case tcell.KeyDown:
		game.snake.MoveDown()
	case tcell.KeyLeft:
		game.snake.MoveLeft()
	}
}

func (game *Game) update() {
	game.snake.ResetMovement()

	head := game.snake.body[0]
	var newHead point

	switch game.snake.direction {
	case Up:
		fallthrough
	case 107:
		newHead = point{head.x, head.y - 1}
	case Down:
		fallthrough
	case 106:
		newHead = point{head.x, head.y + 1}
	case Left:
		fallthrough
	case 108:
		newHead = point{head.x - 1, head.y}
	case Right:
		fallthrough
	case 104:
		newHead = point{head.x + 1, head.y}
	}

	// Check if the snake collides with the walls or itself
	if newHead.x <= 0 || newHead.x >= width || newHead.y <= 0 || newHead.y >= height-1 {
		game.snake.KillIt()
		game.status = Failed

		return
	}

	for _, p := range game.snake.body[1:] {
		if newHead == p {
			game.snake.KillIt()
			game.status = Failed

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

	// game.status =
}

func (game *Game) draw() {
	(*game.Screen).Clear()

	// Draw snake
	for i, p := range game.snake.body {
		char := ' '
		if i == 0 {
			char = '@' // Head

			(*game.Screen).SetContent(p.x, p.y, rune(char), nil, game.headStyle)
		} else {
			char = '#' // Body

			(*game.Screen).SetContent(p.x, p.y, rune(char), nil, game.bodyStyle)
		}
	}

	// Draw food
	(*game.Screen).SetContent(game.food.x, game.food.y, '$', nil, game.foodStyle)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", game.score)
	for i, char := range scoreStr {
		(*game.Screen).SetContent(i, height+1, rune(char), nil, game.textStyle)
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
		(*game.Screen).SetContent(0, y, '|', nil, game.borderStyle)
	}
}

func (game *Game) drawRightBorder() {
	for y := 1; y < height-1; y += 1 {
		(*game.Screen).SetContent(width, y, '|', nil, game.borderStyle)
	}
}

func (game *Game) drawTopBorder() {
	// Top-Left corner
	(*game.Screen).SetContent(0, 0, '/', nil, game.borderStyle)

	for x := 1; x < width; x += 1 {
		(*game.Screen).SetContent(x, 0, '-', nil, game.borderStyle)
	}

	// Top-Right corner
	(*game.Screen).SetContent(width, 0, '\\', nil, game.borderStyle)
}

func (game *Game) drawBottomBorder() {
	// Bottom-Left corner
	(*game.Screen).SetContent(0, height-1, '\\', nil, game.borderStyle)

	for x := 1; x < width; x += 1 {
		(*game.Screen).SetContent(x, height-1, '-', nil, game.borderStyle)
	}

	// Bottom-Right corner
	(*game.Screen).SetContent(width, height-1, '/', nil, game.borderStyle)
}

func (game *Game) placeFood() {
	game.food = point{1 + rand.Intn(width-2), 1 + rand.Intn(height-2)}
}
