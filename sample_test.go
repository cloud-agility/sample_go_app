/* file: $GOPATH/src/godogs/godogs_test.go */
package main

import (
  "fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
  "github.com/DATA-DOG/godog/gherkin"
)

var board *Board
func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func anEmptyBoard() error {
	return nil
}

func iEvolveIt() error {
  board = board.Evolve()
	return nil
}

func aSingleCellOnTheBoard() error {
  board.Set(1,1,true)
	return nil
}

func aXBoardWithTheFollowing(x, y int, data *gherkin.DocString) error {
  board = EmptyBoard(x,y)
  board.Populate(data.Content)
	return nil
}

func itShouldBeEmpty() error {
  empty := EmptyBoard(3,3)
  if !board.Equals(empty) {
          return fmt.Errorf("expected board to be %v, but there is %v", empty, board)
  }
  return nil
}

func itShouldBeLikeTheFollowing(data *gherkin.DocString) error {
  expected := EmptyBoard(board.x,board.y)
  expected.Populate(data.Content)
  if !board.Equals(expected) {
          return fmt.Errorf("expected board to be %v, but there is %v", expected, board)
  }
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^an empty board$`, anEmptyBoard)
	s.Step(`^I evolve it$`, iEvolveIt)
	s.Step(`^it should be empty$`, itShouldBeEmpty)
	s.Step(`^a single cell on the board$`, aSingleCellOnTheBoard)
  s.Step(`^a (\d+) x (\d+) board with the following$`, aXBoardWithTheFollowing)
  s.Step(`^it should be like the following$`, itShouldBeLikeTheFollowing)

	s.BeforeScenario(func(interface{}) {
		board = EmptyBoard(3,3) // clean the state before every scenario
	})
}
