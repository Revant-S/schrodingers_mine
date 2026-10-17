package deduce

import (
	"Schrodingers_Mine/protocol"
	"Schrodingers_Mine/utils/queue"
	"context"
)

// Deduce
// I wanted to replace D here with S .
// But if in case this goes public some day, it will be embarrassing.
// For some reason I want to write some lines for this bot and you can replace D with S everywhere below
// May thy deduction conquer the world

var drow = [8]int{-1, -1, -1, 1, 1, 1, 0, 0}
var dcol = [8]int{-1, 0, 1, -1, 0, 1, -1, 1}

type square struct {
	row int
	col int
}

type Deduce struct {
	safeSquares   *queue.Queue[square]
	rows          int
	cols          int
	mines         int
	minedSquares  map[square]struct{}
	relaxedValues map[square]int
}

func (d *Deduce) NewGame(rows, cols, mines int) {
	d.safeSquares = queue.New[square]()
	d.rows = rows
	d.cols = cols
	d.mines = mines
	d.minedSquares = make(map[square]struct{})
	d.relaxedValues = make(map[square]int)
}

func (d *Deduce) NextMove(ctx context.Context, st *protocol.State) (row, col int, err error) {
	if st.Moves == 0 {
		return 0, 0, nil
	}
	board := st.Board
	d.updateRelaxedSquares(board)

	if sq, ok := d.popCoveredSafe(board); ok {
		return sq.row, sq.col, nil
	}

	numbered, relaxationSquares := d.findMines(board)
	d.relaxAll(board, relaxationSquares)
	d.queueSafeSquares(board, numbered)

	if sq, ok := d.popUnminedSafe(); ok {
		return sq.row, sq.col, nil
	}
	if sq, ok := d.fallbackMove(board); ok {
		return sq.row, sq.col, nil
	}
	return 0, 0, nil // cant reach here ,.. if it does ... God exists
}
