package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/ofer987/snake/models"
)

type GameStatus uint16

const (
	Started GameStatus = iota
	Success
	Failed
)

type Game struct {
	blockedCells [][]bool
	food         models.Point
	score        int
	status       GameStatus
	Screen       *tcell.Screen
	event        Event
	snake        models.Snake
	borderStyle  tcell.Style
	snakeStyle   tcell.Style
	headStyle    tcell.Style
	bodyStyle    tcell.Style
	foodStyle    tcell.Style
	textStyle    tcell.Style
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
	game.borderStyle = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	game.snakeStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	game.headStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlueViolet)
	game.bodyStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
	game.foodStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	game.textStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	game.Screen = screen

	game.snake = models.CreateSnake()

	game.blockedCells = make([][]bool, height)

	game.resetBlockedCells()
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
	if game.snake.GetMovement() == models.Changed {
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

	head := game.snake.Body[0]
	var newHead models.Point

	switch game.snake.GetDirection() {
	case models.Up:
		fallthrough
	case 107:
		newHead = models.Point{X: head.X, Y: head.Y - 1}
	case models.Down:
		fallthrough
	case 106:
		newHead = models.Point{X: head.X, Y: head.Y + 1}
	case models.Left:
		fallthrough
	case 108:
		newHead = models.Point{X: head.X - 1, Y: head.Y}
	case models.Right:
		fallthrough
	case 104:
		newHead = models.Point{X: head.X + 1, Y: head.Y}
	}

	// Check if the snake collides with the walls or itself
	if newHead.X <= 0 || newHead.X >= width || newHead.Y <= 0 || newHead.Y >= height-1 {
		game.snake.KillIt()
		game.status = Failed

		return
	}

	for _, p := range game.snake.Body[1:] {
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
		game.snake.Body = game.snake.Body[:len(game.snake.Body)-1]
	}

	game.snake.Body = append([]models.Point{newHead}, game.snake.Body...)

	// game.status =
}

func (game *Game) draw() {
	game.resetBlockedCells()
	(*game.Screen).Clear()

	// Draw snake
	for i, p := range game.snake.Body {
		char := ' '
		if i == 0 {
			char = '@' // Head

			game.setCell(p.X, p.Y, rune(char), game.headStyle)
		} else {
			char = '#' // Body

			game.setCell(p.X, p.Y, rune(char), game.bodyStyle)
		}
	}

	// Draw food
	game.setCell(game.food.X, game.food.Y, '$', game.foodStyle)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", game.score)
	for i, char := range scoreStr {
		game.setCell(i, height+1, rune(char), game.foodStyle)
	}

	game.drawTopBorder()
	game.drawBottomBorder()
	game.drawLeftBorder()
	game.drawRightBorder()

	if !game.snake.IsAlive() {
		message := "You have lost"
		for i, char := range message {
			game.setCell(i, height+2, rune(char), game.snakeStyle)
		}
	}

	(*game.Screen).Show()
}

func (game *Game) drawLeftBorder() {
	for y := 1; y < height-1; y += 1 {
		game.setCell(0, y, '|', game.borderStyle)
	}
}

func (game *Game) drawRightBorder() {
	for y := 1; y < height-1; y += 1 {
		game.setCell(width, y, '|', game.borderStyle)
	}
}

func (game *Game) drawTopBorder() {
	// Top-Left corner
	game.setCell(0, 0, '/', game.borderStyle)

	for x := 1; x < width; x += 1 {
		game.setCell(x, 0, '–', game.borderStyle)
	}

	// Top-Right corner
	game.setCell(width, 0, '\\', game.borderStyle)
}

func (game *Game) drawBottomBorder() {
	// Bottom-Left corner
	game.setCell(0, height-1, '\\', game.borderStyle)

	for x := 1; x < width; x += 1 {
		game.setCell(x, height-1, '–', game.borderStyle)
	}

	// Bottom-Right corner
	game.setCell(width, height-1, '/', game.borderStyle)
}

func (game *Game) placeFood() {
	for {
		x := 1 + rand.Intn(width-2)
		y := 1 + rand.Intn(height-2)

		if game.blockedCells[x][y] == false {
			game.food = models.Point{X: x, Y: y}

			break
		}
	}
}

func (game *Game) setCell(x int, y int, primary rune, style tcell.Style) {
	game.blockedCells[x][y] = true

	(*game.Screen).SetContent(x, y, primary, nil, style)
}

func (game *Game) resetBlockedCells() {
	game.blockedCells = make([][]bool, width+1)

	for i := range width + 1 {
		game.blockedCells[i] = make([]bool, height+3)
	}
}
