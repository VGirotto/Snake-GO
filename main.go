package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

const sizex = 20
const sizey = 40

type pos struct {
	x, y int
}

type game struct {
	field        [sizex][sizey]int // 0 = empty, 1 = limits h, 2 = limits v, 3 = snake, 4 = food
	pont         int
	snake        []pos
	generateFood bool
	gameOver     bool
	dir          string
}

var waitGo sync.WaitGroup

func main() {

	// termbox to get key input
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)

	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	g := &game{pont: 0, generateFood: true} //inicialize the game
	g.snake = append(g.snake, pos{10, 20})  //inicialize the snake

	g.fieldLimits() //define the walls
	g.dir = "r"     //first direction
	g.gameOver = false

	waitGo.Add(1)

	go g.play() //call the main function

	for g.gameOver == false { //get key input
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'w':
				if g.dir != "d" {
					g.dir = "u"
				}
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 's':
				if g.dir != "u" {
					g.dir = "d"
				}
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'a':
				if g.dir != "r" {
					g.dir = "l"
				}
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'd':
				if g.dir != "l" {
					g.dir = "r"
				}
			case ev.Key == termbox.KeyEsc:
				g.gameOver = true
			}
			time.Sleep(time.Millisecond * 200)
		}
	}

	waitGo.Wait() //wait gameover

	fmt.Println("Game Over!")
	fmt.Println("Score: ", g.pont)
    fmt.Println("Wait 7 seconds")
	time.Sleep(time.Second * 7)
}

func (g *game) play() {

	for g.gameOver == false {
		g.drawSnake()
		g.food()
		g.drawField()
		g.walk()
		g.checkEat()
		time.Sleep(time.Second / 6) //define the velocity of the game

	}
	waitGo.Done()
}

func (g *game) walk() {

	if !(g.snake[len(g.snake)-1].x == 0 && g.snake[len(g.snake)-1].y == 0) { //if snake just eated food, dont delete the last position
		g.field[g.snake[len(g.snake)-1].x][g.snake[len(g.snake)-1].y] = 0 //delete last position
	}

	for i := len(g.snake) - 1; i > 0; i-- {
		g.snake[i] = g.snake[i-1]
	}

	switch g.dir {
	case "r":
		g.checkGameOver(1)
		g.snake[0].y++
	case "l":
		g.checkGameOver(2)
		g.snake[0].y--
	case "d":
		g.checkGameOver(3)
		g.snake[0].x++
	case "u":
		g.checkGameOver(4)
		g.snake[0].x--
	}
}

func (g *game) drawField() { //draw field

	for i := 0; i < sizex; i++ {
		for j := 0; j < sizey; j++ {

			switch g.field[i][j] {
			case 0:
				fmt.Print(" ")
			case 1:
				fmt.Print("—")
			case 2:
				fmt.Print("|")
			case 3:
				fmt.Print("O")
			case 4:
				fmt.Print("x")
			case 5:
				fmt.Print("@")
			}

		}
		fmt.Println()
	}
}

func (g *game) fieldLimits() { //define the limits of the field

	for i, j := 0, 0; i < sizex; i++ { //first column
		g.field[i][j] = 2
	}

	for i, j := 0, sizey-1; i < sizex; i++ { //last column
		g.field[i][j] = 2
	}

	for i, j := 0, 0; j < sizey; j++ { //first line
		g.field[i][j] = 1
	}

	for i, j := sizex-1, 0; j < sizey; j++ { //last line
		g.field[i][j] = 1
	}
}

func (g *game) drawSnake() { //define the pos of the snake in the field
	for i := 0; i < len(g.snake); i++ {
		if !(g.snake[i].x == 0) {
			g.field[g.snake[i].x][g.snake[i].y] = 3
		}
	}
}

func (g *game) food() {
	if g.generateFood {

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		x := r1.Intn(sizex-2) + 1
		y := r1.Intn(sizey-2) + 1

		for g.field[x][y] == 3 {
			x = r1.Intn(sizex-2) + 1
			y = r1.Intn(sizey-2) + 1
		}
		g.field[x][y] = 4
		g.generateFood = false

	}
}

func (g *game) checkEat() {
	if g.field[g.snake[0].x][g.snake[0].y] == 4 { //if head of the snake is in same position of the food
		g.pont++
		g.generateFood = true
		g.snake = append(g.snake, pos{0, 0})
	}
}

func (g *game) checkGameOver(num int) {
	switch num {
	case 1: //right
		if g.field[g.snake[0].x][g.snake[0].y+1] == 1 || g.field[g.snake[0].x][g.snake[0].y+1] == 2 || g.field[g.snake[0].x][g.snake[0].y+1] == 3 {
			//gameover
			g.gameOver = true
			g.snake[0].y++
		}

	case 2: //left
		if g.field[g.snake[0].x][g.snake[0].y-1] == 1 || g.field[g.snake[0].x][g.snake[0].y-1] == 2 || g.field[g.snake[0].x][g.snake[0].y-1] == 3 {
			//gameover
			g.gameOver = true
			g.snake[0].y--
		}

	case 3: //down
		if g.field[g.snake[0].x+1][g.snake[0].y] == 1 || g.field[g.snake[0].x+1][g.snake[0].y] == 2 || g.field[g.snake[0].x+1][g.snake[0].y] == 3 {
			//gameover
			g.gameOver = true
			g.snake[0].x++
		}

	case 4: //up
		if g.field[g.snake[0].x-1][g.snake[0].y] == 1 || g.field[g.snake[0].x-1][g.snake[0].y] == 2 || g.field[g.snake[0].x-1][g.snake[0].y] == 3 {
			//gameover
			g.gameOver = true
			g.snake[0].x--
		}

	}

	if g.gameOver == true {
		g.drawSnake()
		g.field[g.snake[0].x][g.snake[0].y] = 5
		g.drawField()
		g.snake[0].x = 1
		g.snake[0].y = 1
	}
}
