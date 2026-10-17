package deduce

import "Schrodingers_Mine/protocol"

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

func (d *Deduce) updateRelaxedSquares(board [][]int) {
	for sq, val := range d.relaxedValues {
		board[sq.row][sq.col] = val
	}
}
