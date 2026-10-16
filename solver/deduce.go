package solver

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
	safeSquares  *queue.Queue[square]
	rows         int
	cols         int
	mines        int
	minedSquares map[square]struct{}
}

func (d *Deduce) NewGame(rows, cols, mines int) {
	d.safeSquares = queue.New[square]()
	d.rows = rows
	d.cols = cols
	d.mines = mines
	d.minedSquares = make(map[square]struct{})
}

func (d *Deduce) markMined(sq square) {
	d.minedSquares[sq] = struct{}{}
}

func (d *Deduce) checkCovered(board [][]int, sq square) bool {
	return board[sq.row][sq.col] == protocol.Covered
}

func (d *Deduce) containsNumb(board [][]int, sq square) bool {
	v := board[sq.row][sq.col]
	return v >= 1 && v <= 8
}

func (d *Deduce) isValidCoordinate(sq square) bool {
	return sq.row >= 0 && sq.col >= 0 && sq.row < d.rows && sq.col < d.cols
}

func (d *Deduce) numOfCoveredNeighbours(board [][]int, sq square) int {
	count := 0
	for idx := 0; idx < 8; idx++ {
		nb := square{row: sq.row + drow[idx], col: sq.col + dcol[idx]}
		if d.isValidCoordinate(nb) && d.checkCovered(board, nb) {
			count++
		}
	}
	return count
}

func (d *Deduce) neighbours(board [][]int, sq square) []square {
	var list []square
	for idx := 0; idx < 8; idx++ {
		nb := square{row: sq.row + drow[idx], col: sq.col + dcol[idx]}
		if d.isValidCoordinate(nb) && d.checkCovered(board, nb) {
			list = append(list, nb)
		}

	}
	return list
}

func (d *Deduce) relaxNeighbours(board [][]int, sq square) {
	for idx := 0; idx < 8; idx++ {
		nb := square{row: sq.row + drow[idx], col: sq.col + dcol[idx]}
		if d.isValidCoordinate(nb) && d.containsNumb(board, nb) {
			board[nb.row][nb.col] -= 1
		}
	}
}

func (d *Deduce) NextMove(ctx context.Context, st *protocol.State) (row, col int, err error) {
	if st.Moves == 0 {
		return 0, 0, nil
	}
	board := st.Board
	// if safe squares are known, and they are still available play them 1st
	for {
		front, ok := d.safeSquares.Pop()
		if !ok {
			break
		}
		if d.checkCovered(board, front) {
			return front.row, front.col, nil
		}
	}
	numbered := make([]square, 0, d.rows*d.cols)
	// go through the whole board and check for the squares where the number is equal to
	// number of its covered inputs
	for r := range board {
		for c := range board[r] {
			sq := square{row: r, col: c}
			if !d.containsNumb(board, sq) {
				continue
			}
			numbered = append(numbered, sq)
			// if it contains a number and the number is same as the number of covered neighbours
			// All those neighbours are mined
			if board[sq.row][sq.col] == d.numOfCoveredNeighbours(board, sq) {
				nbrs := d.neighbours(board, sq)
				//mined := make([]bool, len(nbrs))
				for _, nb := range nbrs {
					if _, ok := d.minedSquares[nb]; ok {
						continue
					}
					d.markMined(nb)
					d.relaxNeighbours(board, nb)
				}
			}
		}
	}
	// Now check for the safe squares
	for _, s := range numbered {
		// if a numbered square is now 0 . all its neighbours are safe now
		if board[s.row][s.col] == 0 {
			nbrs := d.neighbours(board, s)
			for _, nbr := range nbrs {
				if d.checkCovered(board, nbr) {
					d.safeSquares.Push(nbr)
				}
			}
		}
	}
	// return a safe move
	if !d.safeSquares.IsEmpty() {
		front, _ := d.safeSquares.Pop()
		return front.row, front.col, nil
	}
	//d.safeSquares.Print()
	return 0, 0, nil
}
