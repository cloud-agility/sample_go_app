// A refactored implementation of Conway's Game of Life.
package main

import (
	"bytes"
  "strings"
	"fmt"
	"math/rand"
  "net/http"
)

var CLEARSCREEN string = "\033[H\033[2J"

type Board struct {
	cells [][]bool
	x, y int
}

func EmptyBoard(x, y int) *Board {
	cells := make([][]bool, y)
	for i := range cells {
		cells[i] = make([]bool, x)
	}
	return &Board{cells: cells, x: x, y: y}
}

func (board *Board) RandomizePopulation(ratio int) {
	for i := 0; i < (board.x * board.y / ratio); i++ {
		board.Set(rand.Intn(board.x), rand.Intn(board.y), true)
	}
}

func (board *Board) Set(x, y int, state bool) {
	board.cells[y][x] = state
}

func (board *Board) IsAliveAt(x, y int) bool {
  // If the x or y coordinates are outside the field boundaries they are wrapped
  // toroidally. For instance, an x value of -1 is treated as width-1.
  x += board.x
  x %= board.x
  y += board.y
  y %= board.y
	return board.cells[y][x]
}

func (board *Board) EvolveCell(x, y int) bool {
	neighbourCount := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if (NotThisCell(dx,dy)) && board.IsAliveAt(x+dx, y+dy) {
				neighbourCount++
			}
		}
	}
	return neighbourCount == 2 && board.IsAliveAt(x, y) || neighbourCount == 3
}

func NotThisCell(x, y int) bool {
  return (x != 0 || y != 0 )
}

func (board *Board) Evolve() *Board {
  nextgen := EmptyBoard(board.x, board.y)
	for y := 0; y < board.y; y++ {
		for x := 0; x < board.x; x++ {
			nextgen.Set(x, y, board.EvolveCell(x, y))
		}
	}
	return &Board{cells: nextgen.cells, x: nextgen.x, y: nextgen.y}
}

// called implicitly by Print()
func (board *Board) String() string {
	var buffer bytes.Buffer
	for y := 0; y < board.y; y++ {
		for x := 0; x < board.x; x++ {
			b := byte('.')
			if board.IsAliveAt(x, y) {
				b = '*'
			}
			buffer.WriteByte(b)
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

// ServeHTTP is called when board is used as the Handler for requests
func (board *Board) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, board)
  board = board.Evolve()
}


func (board *Board) Populate(data string) {
  bytes := []byte(strings.Replace(data, "\n", "", -1))

  for y := 0; y < board.y; y++ {
    for x := 0; x < board.x; x++ {
      if bytes[y*board.y+x]=='*' {
        board.Set(x,y,true)
      }
    }
  }
}

// Equals provides equality testing
func (board *Board) Equals(test *Board) bool {
  return board.x == test.x && board.y == test.y && board.String() == test.String()
}

func main() {
  board := EmptyBoard(15,15)
  board.RandomizePopulation(4)
  http.Handle("/", board)
  http.ListenAndServe(":8080", nil)
}
