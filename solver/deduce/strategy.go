package deduce

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
