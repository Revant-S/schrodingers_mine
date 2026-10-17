package deduce

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
