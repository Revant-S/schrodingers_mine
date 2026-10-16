package solver

import (
	"Schrodingers_Mine/protocol"
	"Schrodingers_Mine/utils/queue"
	"context"
	"fmt"
	"strings"
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
		_, mined := d.minedSquares[nb]
		if d.isValidCoordinate(nb) && d.containsNumb(board, nb) && !mined {
			board[nb.row][nb.col] -= 1
			d.relaxedValues[nb] = board[nb.row][nb.col]
		}
	}
}

func (d *Deduce) updateRelaxedSquares(board [][]int) {
	for sq, val := range d.relaxedValues {
		board[sq.row][sq.col] = val
	}
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

// popCoveredSafe pops already calculated safe squares until it finds one that is still covered.
func (d *Deduce) popCoveredSafe(board [][]int) (square, bool) {
	for {
		front, ok := d.safeSquares.Pop()
		if !ok {
			return square{}, false
		}
		if d.checkCovered(board, front) {
			return front, true
		}
	}
}

// findMines collects the numbered squares and marks as mined the covered neighbours of every
// number that equals its count of covered neighbours. It returns the numbered squares and the
// newly mined squares, whose neighbours still need relaxing.
func (d *Deduce) findMines(board [][]int) (numbered, relaxationSquares []square) {
	numbered = make([]square, 0, d.rows*d.cols)
	relaxationSquares = make([]square, 0, d.rows*d.cols)
	for r := range board {
		for c := range board[r] {
			sq := square{row: r, col: c}
			if !d.containsNumb(board, sq) {
				continue
			}
			numbered = append(numbered, sq)
			if board[sq.row][sq.col] == d.numOfCoveredNeighbours(board, sq) {
				nbrs := d.neighbours(board, sq)
				for _, nb := range nbrs {
					if _, ok := d.minedSquares[nb]; ok {
						continue
					}
					d.markMined(nb)
					relaxationSquares = append(relaxationSquares, nb)
				}
			}
		}
	}
	return numbered, relaxationSquares
}

// relaxAll relaxes the neighbours of every given mined square.
func (d *Deduce) relaxAll(board [][]int, mined []square) {
	for _, relaxationSquare := range mined {
		d.relaxNeighbours(board, relaxationSquare)
	}
}

// queueSafeSquares pushes the covered, unmined neighbours of every numbered square that is now 0.
func (d *Deduce) queueSafeSquares(board [][]int, numbered []square) {
	for _, s := range numbered {
		if board[s.row][s.col] == 0 {
			nbrs := d.neighbours(board, s)
			for _, nbr := range nbrs {
				_, mined := d.minedSquares[nbr]
				if d.checkCovered(board, nbr) && !mined {
					d.safeSquares.Push(nbr)
				}
			}
		}
	}
}

// popUnminedSafe pops safe squares until it finds one that is not marked as mined.
func (d *Deduce) popUnminedSafe() (square, bool) {
	for {
		if !d.safeSquares.IsEmpty() {
			front, _ := d.safeSquares.Pop()
			if _, mined := d.minedSquares[front]; !mined {
				return front, true
			}
		} else {
			return square{}, false
		}
	}
}

// fallbackMove returns the first square in reading order that is not a number, not 0 and not mined.
func (d *Deduce) fallbackMove(board [][]int) (square, bool) {
	for i := 0; i < d.rows; i++ {
		for j := 0; j < d.cols; j++ {
			sq := square{row: i, col: j}
			_, mined := d.minedSquares[sq]
			if !d.containsNumb(board, sq) && board[i][j] != 0 && !mined {
				return sq, true
			}
		}
	}
	return square{}, false
}

// printBoard shows the board as the solver sees it.
//
//	.   uncovered, no mines around (0)
//	1-8 number (after any relaxing you did)
//	#   covered, unknown
//	M   covered, marked as a mine by the solver
//	X   mine shown by the board after a loss (-3)
func (d *Deduce) printBoard(board [][]int) {
	var b strings.Builder

	b.WriteString("    ")
	for c := range board[0] {
		fmt.Fprintf(&b, "%2d ", c)
	}
	b.WriteString("\n    " + strings.Repeat("---", len(board[0])) + "\n")

	for r, row := range board {
		_, err := fmt.Fprintf(&b, "%2d |", r)
		if err != nil {
			return
		}

		for c, v := range row {
			var ch string
			switch {
			case v == protocol.Covered:
				if _, ok := d.minedSquares[square{row: r, col: c}]; ok {
					ch = "M"
				} else {
					ch = "#"
				}
			case v == protocol.Mine:
				ch = "X"
			case v == 0:
				ch = "."
			default:
				ch = fmt.Sprint(v)
			}
			_, err := fmt.Fprintf(&b, "%2s ", ch)
			if err != nil {
				return

			}
		}
		b.WriteString("\n")
	}
	fmt.Print(b.String())
}
