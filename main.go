package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const boardY int = 22
const boardX int = 12

type Point struct {
	y int
	x int
}

type shape interface {
	init()
	move(y int, x int)
	getLocation() []Point
	getColor() tcell.Color
}

type shapeStruct struct {
	location []Point
	color    tcell.Color
}

type line struct {
	shapeStruct
}

type plus struct {
	shapeStruct
}

type square struct {
	shapeStruct
}

func (s *square) init() {
	s.location = []Point{{1, 1}, {1, 2}, {2, 1}, {2, 2}}
	s.color = tcell.ColorBlue
}

func (p *plus) init() {
	p.location = []Point{{1, 1}, {1, 2}, {1, 3}, {2, 2}}
	p.color = tcell.ColorYellow
}

func (l *line) init() {
	l.location = []Point{{1, 1}, {1, 2}, {1, 3}, {1, 4}}
	l.color = tcell.ColorPurple
}

func (s *shapeStruct) move(y int, x int) {
	for _, l := range s.location {
		if board[l.y][l.x+x] != tcell.ColorWhite {
			return
		}
		if board[l.y+y][l.x] != tcell.ColorWhite {
			stopMove()
			return
		}
	}
	for i, l := range s.location {
		l.y += y
		l.x += x
		s.location[i] = l
	}
}

func (s *shapeStruct) getColor() tcell.Color {
	return s.color
}

func (s *shapeStruct) getLocation() []Point {
	return s.location
}

var board [boardY][boardX]tcell.Color

var (
	table   *tview.Table
	app     *tview.Application
	current shape
)

func stopMove() {
	color := current.getColor()
	for _, s := range current.getLocation() {
		board[s.y][s.x] = color
	}
	r := rand.Intn(3)
	switch r {
	case 0:
		current = &line{}
		break
	case 1:
		current = &plus{}
	case 2:
		current = &square{}
	}
	current.init()
}

func updateTime() {
	for {
		time.Sleep(1 * time.Second)
		app.QueueUpdateDraw(func() {
			current.move(1, 0)
			drawBoard()
		})
	}
}

func drawBoard() {
	for y, l := range board {
		for x, v := range l {
			// b := table.GetCell(y, 2*x).BackgroundColor
			// fmt.Println(b)
			//if table.GetCell(y, 2*x).BackgroundColor != v {
			table.GetCell(y, 2*x).SetBackgroundColor(v)
			table.GetCell(y, 2*x+1).SetBackgroundColor(v)
			// table.SetCell(y, 2*x, tview.NewTableCell("").SetBackgroundColor(v))
			// table.SetCell(y, 2*x+1, tview.NewTableCell("").SetBackgroundColor(v))
			//}
		}
	}
	for _, s := range current.getLocation() {
		table.GetCell(s.y, 2*s.x).SetBackgroundColor(current.getColor())
		table.GetCell(s.y, 2*s.x+1).SetBackgroundColor(current.getColor())
		// table.SetCell(s.y, 2*s.x, tview.NewTableCell("").SetBackgroundColor(current.getColor()))
		// table.SetCell(s.y, 2*s.x+1, tview.NewTableCell("").SetBackgroundColor(current.getColor()))
	}
}

func initBoard() {
	for y := 0; y < boardY; y++ {
		for x := 0; x < boardX; x++ {
			if x == 0 || x == boardX-1 || y == 0 || y == boardY-1 {
				board[y][x] = tcell.ColorBlack
			} else {
				board[y][x] = tcell.ColorWhite
			}
			table.SetCell(y, 2*x, tview.NewTableCell("").SetBackgroundColor(board[y][x]))
			table.SetCell(y, 2*x+1, tview.NewTableCell("").SetBackgroundColor(board[y][x]))
		}
	}
}

func checkBoard() {
	//for
}

func main() {
	app = tview.NewApplication()
	table = tview.NewTable()
	initBoard()
	current = &line{}
	current.init()
	table.SetBorder(true).SetTitle("Hello, world!")
	drawBoard()
	flex := tview.NewFlex().
		AddItem(table, 40, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			current.move(0, -1)
			break
		case tcell.KeyRight:
			current.move(0, 1)
			break
		case tcell.KeyDown:
			current.move(1, 0)
			break
		}
		drawBoard()
		return event
	})
	go updateTime()
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
