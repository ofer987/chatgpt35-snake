package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width  = 30
	height = 20
)

type point struct {
	x, y int
}

type snake struct {
	body      []point
	direction termbox.Key
	alive     bool
}

var (
	food     point
	score    int
	gameOver bool
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	snake := snake{
		body:      []point{{5, 5}, {4, 5}, {3, 5}},
		direction: termbox.KeyArrowRight,
		alive:     true,
	}

	placeFood()

	gameLoop(&snake)

	// if !gameOver {
	// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// 	termbox.SetCell(width/2-5, height/2, 'G', termbox.ColorRed, termbox.ColorDefault)
	// 	// message := "You have lost"
	// 	//
	// 	// for i, ch := range message {
	// 	// 	termbox.SetCell(i+1, 5, ch, termbox.ColorGreen, termbox.ColorDefault)
	// 	// }
	// 	termbox.Flush()
	// 	time.Sleep(time.Second * 2)
	// }
}

func gameLoop(s *snake) {
	ticker := time.Tick(100 * time.Millisecond)

	for {
		select {
		case ev := <-eventQueue():
			if s.alive {
				handleEventWhenAlive(ev, s)
			} else {
				handleEventWhenDead(ev, s)
			}
		case <-ticker:
			if gameOver {
				return
			}
			update(s)
			draw(s)
		}
	}
}

func eventQueue() <-chan termbox.Event {
	ch := make(chan termbox.Event)
	go func() {
		for {
			ch <- termbox.PollEvent()
		}
	}()

	return ch
}

func handleEventWhenAlive(ev termbox.Event, s *snake) {
	if ev.Type == termbox.EventKey && ev.Ch == 0 {
		// termbox.SetCell(1, height+4, rune(ev.Key), termbox.ColorGreen, termbox.ColorDefault)
		switch ev.Key {
		case termbox.KeyArrowUp:
			fallthrough
		case 107:
			if s.direction != termbox.KeyArrowDown {
				s.direction = ev.Key
			}
		case termbox.KeyArrowDown:
			fallthrough
		case 106:
			if s.direction != termbox.KeyArrowUp {
				s.direction = ev.Key
			}
		case termbox.KeyArrowLeft:
			fallthrough
		case 108:
			if s.direction != termbox.KeyArrowRight {
				s.direction = ev.Key
			}
		case termbox.KeyArrowRight:
			fallthrough
		case 104:
			if s.direction != termbox.KeyArrowLeft {
				s.direction = ev.Key
				// fmt.Printf("%d", ev.Key)
			}
		}
	}
}

func handleEventWhenDead(ev termbox.Event, s *snake) {
	if ev.Type == termbox.EventKey && ev.Ch == 0 {
		switch ev.Key {
		case termbox.KeyEnter:
			fallthrough
		case termbox.KeyEsc:
			if !s.alive {
				gameOver = true
				return
			}
		}
	}
}

func update(s *snake) {
	head := s.body[0]
	var newHead point

	switch s.direction {
	case termbox.KeyArrowUp:
		fallthrough
	case 107:
		newHead = point{head.x, head.y - 1}
	case termbox.KeyArrowDown:
		fallthrough
	case 106:
		newHead = point{head.x, head.y + 1}
	case termbox.KeyArrowLeft:
		fallthrough
	case 108:
		newHead = point{head.x - 1, head.y}
	case termbox.KeyArrowRight:
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

func draw(s *snake) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw snake
	for i, p := range s.body {
		char := ' '
		if i == 0 {
			char = '@' // Head
		} else {
			char = '#' // Body
		}

		termbox.SetCell(p.x, p.y, char, termbox.ColorGreen, termbox.ColorDefault)
	}

	// Draw food
	termbox.SetCell(food.x, food.y, '$', termbox.ColorYellow, termbox.ColorDefault)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", score)
	for i, ch := range scoreStr {
		termbox.SetCell(i, height+5, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	drawTopBorder()
	drawBottomBorder()
	drawLeftBorder()
	drawRightBorder()

	if !s.alive {
		message := "You have lost"

		for i, ch := range message {
			termbox.SetCell(i+1, 5, ch, termbox.ColorGreen, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func drawLeftBorder() {
	for y := 1; y < height-1; y += 1 {
		termbox.SetCell(0, y, '|', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func drawRightBorder() {
	for y := 1; y < height-1; y += 1 {
		termbox.SetCell(width, y, '|', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func drawTopBorder() {
	// Top-Left corner
	termbox.SetCell(0, 0, '/', termbox.ColorWhite, termbox.ColorDefault)

	for x := 1; x < width; x += 1 {
		termbox.SetCell(x, 0, '–', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Top-Right corner
	termbox.SetCell(width, 0, '\\', termbox.ColorWhite, termbox.ColorDefault)
}

func drawBottomBorder() {
	// Bottom-Left corner
	termbox.SetCell(0, height-1, '\\', termbox.ColorWhite, termbox.ColorDefault)

	for x := 1; x < width; x += 1 {
		termbox.SetCell(x, height-1, '–', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Bottom-Right corner
	termbox.SetCell(width, height-1, '/', termbox.ColorWhite, termbox.ColorDefault)
}

func placeFood() {
	// rand.Seed(time.Now().UnixNano())
	food = point{1 + rand.Intn(width-2), 1 + rand.Intn(height-2)}
}
