package deduce

import (
	"Schrodingers_Mine/protocol"
	"fmt"
	"strings"
)

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
		_, err := fmt.Fprintf(&b, "%2d ", c)
		if err != nil {
			return
		}
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
