// A refactored implementation of Conway's Game of Life.
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

// Board represents the game of life state on an x by y sized grid
type Board struct {
	cells [][]bool
	x, y  int
}

// EmptyBoard creates an empty grid of size x by y
func EmptyBoard(x, y int) *Board {
	cells := make([][]bool, y)
	for i := range cells {
		cells[i] = make([]bool, x)
	}
	return &Board{cells: cells, x: x, y: y}
}

// RandomizePopulation a portion of the board by the specified ratio
func (board *Board) RandomizePopulation(ratio int) {
	for i := 0; i < (board.x * board.y / ratio); i++ {
		board.Set(rand.Intn(board.x), rand.Intn(board.y), true)
	}
}

// Set a specific grid cell on the board to the specified state
func (board *Board) Set(x, y int, state bool) {
	board.cells[y][x] = state
}

// IsAliveAt returns true if a specified location on the board is alive
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (board *Board) IsAliveAt(x, y int) bool {

	x += board.x
	x %= board.x
	y += board.y
	y %= board.y
	return board.cells[y][x]
}

// EvolveCell returns the next cell state for a given cell location on the board
func (board *Board) EvolveCell(x, y int) bool {
	neighbourCount := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if (NotThisCell(dx, dy)) && board.IsAliveAt(x+dx, y+dy) {
				neighbourCount++
			}
		}
	}
	return neighbourCount == 2 && board.IsAliveAt(x, y) || neighbourCount == 3
}

// NotThisCell returns true if the specified location is anywhere but the origin
func NotThisCell(x, y int) bool {
	return (x != 0 || y != 0)
}

// Evolve mutates the state of a board to the next evolution of the game world
func (board *Board) Evolve() {
	nextgen := EmptyBoard(board.x, board.y)
	for y := 0; y < board.y; y++ {
		for x := 0; x < board.x; x++ {
			nextgen.Set(x, y, board.EvolveCell(x, y))
		}
	}
	board.x = nextgen.x
	board.y = nextgen.y
	board.cells = nextgen.cells
}

// String is called implicitly by Print() and renders the board state as text
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
// It has the quirky approach of also evolving the board...which is a terrible idea in general
func (board *Board) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	board.Evolve()
	fmt.Fprint(w, board)
}

// Populate the board with a given text input representation of a game world
func (board *Board) Populate(data string) {
	bytes := []byte(strings.Replace(data, "\n", "", -1))

	for y := 0; y < board.y; y++ {
		for x := 0; x < board.x; x++ {
			if bytes[y*board.y+x] == '*' {
				board.Set(x, y, true)
			}
		}
	}
}

// Equals provides equality testing of two boards
func (board *Board) Equals(test *Board) bool {
	return board.x == test.x && board.y == test.y && board.String() == test.String()
}

func main() {
	board := EmptyBoard(15, 15)
	board.RandomizePopulation(4)
	http.Handle("/", board)
	http.ListenAndServe(":8080", nil)
}
