package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const boardX int = 22
const boardY int = 12

type Point struct {
	x int
	y int
}

type shape interface {
	init()
	move(x int, y int)
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

func (l *line) init() {
	l.location = []Point{{1, 1}, {1, 2}, {1, 3}, {1, 4}}
	l.color = tcell.ColorPurple
}

func (s *shapeStruct) move(x int, y int) {
	movePiece(x, y, s.location[:], s.getColor())
}

func (s *shapeStruct) getColor() tcell.Color {
	return s.color
}

func (s *shapeStruct) getLocation() []Point {
	return s.location
}

var board [boardX][boardY]tcell.Color

var (
	table   *tview.Table
	app     *tview.Application
	current shape
)

func stopMove() {
	// for _, s := range current {
	// 	board[s.y][s.x] = "b"
	// }
	// current[0] = Point{1, 1}
	// current[1] = Point{1, 2}
	// current[2] = Point{2, 2}
	// current[3] = Point{2, 3}
}

func drawPiece(location []Point, pieceColor tcell.Color) {

}

func movePiece(x int, y int, location []Point, pieceColor tcell.Color) {
	for _, s := range location {
		if board[s.x][s.y+y] != tcell.ColorWhite {
			return
		}
		if board[s.x+x][s.y] != tcell.ColorWhite {
			stopMove()
			return
		}
	}
	for i, s := range location {
		s.x += x
		s.y += y
		location[i] = s
	}
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
	for x, l := range board {
		for y, v := range l {
			table.SetCell(x, 2*y, tview.NewTableCell("").SetBackgroundColor(v))
			table.SetCell(x, 2*y+1, tview.NewTableCell("").SetBackgroundColor(v))
		}
	}
	for _, s := range current.getLocation() {
		table.SetCell(s.x, 2*s.y, tview.NewTableCell("").SetBackgroundColor(current.getColor()))
		table.SetCell(s.x, 2*s.y+1, tview.NewTableCell("").SetBackgroundColor(current.getColor()))
	}
}

func initBoard() {
	for x := 0; x < boardX; x++ {
		for y := 0; y < boardY; y++ {
			if x == 0 || x == boardX-1 || y == 0 || y == boardY-1 {
				board[x][y] = tcell.ColorBlack
			} else {
				board[x][y] = tcell.ColorWhite
			}
		}
	}
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
			// movePiece(0, -1)
			break
		case tcell.KeyRight:
			current.move(0, 1)
			// movePiece(0, 1)
			break
		case tcell.KeyDown:
			current.move(1, 0)
			// movePiece(1, 0)
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
