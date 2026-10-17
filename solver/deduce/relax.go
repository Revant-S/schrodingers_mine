package deduce

func (d *Deduce) relaxNeighbours(board [][]int, sq square) {
	for idx := 0; idx < 8; idx++ {
		nb := square{row: sq.row + drow[idx], col: sq.col + dcol[idx]}
		_, mined := d.minedSquares[nb]
		if d.isValidCoordinate(nb) && d.containsNumb(board, nb) && !mined && !d.checkCovered(board, nb) {
			board[nb.row][nb.col] -= 1
			d.relaxedValues[nb] = board[nb.row][nb.col]
		}
	}
}

func (d *Deduce) relaxAll(board [][]int, mined []square) {
	for _, relaxationSquare := range mined {
		d.relaxNeighbours(board, relaxationSquare)
	}
}
