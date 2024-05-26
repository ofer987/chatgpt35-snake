package main

type Snake struct {
	body      []point
	direction Direction
	alive     bool
	movement  Movement
}

type Direction uint8

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Movement uint16

const (
	Moving Movement = iota
	Changed
)

func CreateSnake() Snake {
	snake := Snake{
		body:      []point{{5, 5}, {4, 5}, {3, 5}},
		direction: Right,
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
	case Up:
		break
	case Right:
		snake.direction = Up
	case Down:
		break
	case Left:
		snake.direction = Up
	}
}

func (snake *Snake) MoveRight() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case Up:
		snake.direction = Right
	case Right:
		break
	case Down:
		snake.direction = Right
	case Left:
		break
	}
}

func (snake *Snake) MoveDown() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case Up:
		break
	case Right:
		snake.direction = Down
	case Down:
		break
	case Left:
		snake.direction = Down
	}
}

func (snake *Snake) MoveLeft() {
	if snake.GetMovement() == Changed {
		return
	}

	switch snake.direction {
	case Up:
		snake.direction = Left
	case Right:
		break
	case Down:
		snake.direction = Left
	case Left:
		break
	}
}
