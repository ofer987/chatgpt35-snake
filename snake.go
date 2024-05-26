package main

import (
	"github.com/gdamore/tcell/v2"
)

type Snake struct {
	body      []point
	direction tcell.Key
	alive     bool
	movement  Movement
}

type Movement uint16

const (
	Moving Movement = iota
	Changed
)

func CreateSnake() Snake {
	snake := Snake{
		body:      []point{{5, 5}, {4, 5}, {3, 5}},
		direction: tcell.KeyRight,
		alive:     true,
	}

	return snake
}

func (snake *Snake) GetMovement() Movement {
	return snake.movement
}

func (snake *Snake) IsAlive() bool {
	return snake.alive
}

func (snake *Snake) KillIt() {
	snake.alive = false
}

func (snake *Snake) ResetMovement() {
	snake.movement = Moving
}

func (snake *Snake) MoveUp() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case tcell.KeyUp:
		break
	case tcell.KeyRight:
		snake.direction = tcell.KeyUp
	case tcell.KeyDown:
		break
	case tcell.KeyLeft:
		snake.direction = tcell.KeyUp
	}
}

func (snake *Snake) MoveRight() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case tcell.KeyUp:
		snake.direction = tcell.KeyRight
	case tcell.KeyRight:
		break
	case tcell.KeyDown:
		snake.direction = tcell.KeyRight
	case tcell.KeyLeft:
		break
	}
}

func (snake *Snake) MoveDown() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case tcell.KeyUp:
		break
	case tcell.KeyRight:
		snake.direction = tcell.KeyDown
	case tcell.KeyDown:
		break
	case tcell.KeyLeft:
		snake.direction = tcell.KeyDown
	}
}

func (snake *Snake) MoveLeft() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case tcell.KeyUp:
		snake.direction = tcell.KeyLeft
	case tcell.KeyRight:
		break
	case tcell.KeyDown:
		snake.direction = tcell.KeyLeft
	case tcell.KeyLeft:
		break
	}
}
