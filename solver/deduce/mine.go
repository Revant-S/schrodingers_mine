package deduce

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

func (d *Deduce) markMined(sq square) {
	d.minedSquares[sq] = struct{}{}
}
