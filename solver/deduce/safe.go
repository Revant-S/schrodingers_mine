package deduce

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

// relaxAll relaxes the neighbours of every given mined square.

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
