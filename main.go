package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/mattn/go-tty"
)

var score = 0

type Point struct {
	X, Y int
}

func (p1 *Point) isEqual(p2 *Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

var directionMap = map[rune]Point{
	'w': {0, -1},
	's': {0, 1},
	'a': {-1, 0},
	'd': {1, 0},
}

var directionHeadCharMap = map[Point]string{
	{0, -1}: "▲",
	{0, 1}:  "▼",
	{-1, 0}: "◀",
	{1, 0}:  "▶",
}

type Grid struct {
	Width, Height int
	Snake         *Snake
	Food          *Point
	Quit          chan struct{}
}

func NewGrid(width, height int, initialSnakeDir Point) *Grid {
	return &Grid{
		Width:  width,
		Height: height,
		Snake:  NewSnake(directionMap['d']),
		Food:   spawnFood(width, height),
		Quit:   make(chan struct{}),
	}
}

type Snake struct {
	body     []Point
	dir      Point
	headChar string
}

func (s *Snake) changeDir(char rune) {
	switch char {
	case 'w':
		if s.dir.Y == 0 {
			s.dir = directionMap[char]
			s.headChar = directionHeadCharMap[s.dir]
		}
	case 's':
		if s.dir.Y == 0 {
			s.dir = directionMap[char]
			s.headChar = directionHeadCharMap[s.dir]
		}
	case 'a':
		if s.dir.X == 0 {
			s.dir = directionMap[char]
			s.headChar = directionHeadCharMap[s.dir]
		}
	case 'd':
		if s.dir.X == 0 {
			s.dir = directionMap[char]
			s.headChar = directionHeadCharMap[s.dir]
		}
	}
}

func NewSnake(initialDir Point) *Snake {
	return &Snake{
		body:     []Point{{5, 5}},
		dir:      initialDir,
		headChar: directionHeadCharMap[initialDir],
	}
}

func spawnFood(width, height int) *Point {
	return &Point{rand.Intn(width-2) + 1, rand.Intn(height-2) + 1}
}

func initializeGame() *Grid {
	return NewGrid(100, 20, Point{1, 0})
}

func handleInput(grid *Grid) {
	tty, err := tty.Open()
	if err != nil {
		fmt.Println("Error initializing input:", err)
		close(grid.Quit)
		return
	}
	defer tty.Close()

	for {
		char, err := tty.ReadRune()
		if err != nil {
			continue
		}
		if char == 'q' {
			close(grid.Quit)
			return
		} else {
			grid.Snake.changeDir(char)
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	grid := initializeGame()
	go handleInput(grid)
	gameLoop(grid)
}

func gameLoop(grid *Grid) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-grid.Quit:
			return
		case <-ticker.C:
			updateGame(grid)
			renderGame(grid)
		}
	}
}

func updateGame(grid *Grid) {
	snake := grid.Snake
	head := Point{snake.body[0].X + snake.dir.X, snake.body[0].Y + snake.dir.Y}
	if head.X <= 0 || head.X >= grid.Width-1 || head.Y <= 0 || head.Y >= grid.Height-1 {
		close(grid.Quit)
		return
	}

	for _, p := range snake.body {
		if p == head {
			close(grid.Quit)
			return
		}
	}

	snake.body = append([]Point{head}, snake.body...)
	if head == *grid.Food {
		grid.Food = spawnFood(grid.Width, grid.Height)
		score = score + 1
	} else {
		snake.body = snake.body[:len(snake.body)-1]
	}
}

func renderGame(grid *Grid) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for y := range grid.Height {
		for x := range grid.Width {
			if y == 0 || y == grid.Height-1 || x == 0 || x == grid.Width-1 {
				fmt.Print("█")
			} else {
				currPoint := &Point{x, y}
				isSnake := false
				for i, p := range grid.Snake.body {
					if p.isEqual(currPoint) {
						char := "■"
						if i == 0 {
							char = grid.Snake.headChar
						}
						fmt.Print(char)
						isSnake = true
						break
					}
				}
				if !isSnake {
					if grid.Food.isEqual(currPoint) {
						fmt.Print("●")
					} else {
						fmt.Print(" ")
					}
				}
			}
		}
		fmt.Println()
	}
	fmt.Printf("Score: %d", score)
}
