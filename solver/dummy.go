// Package solver is where the actual Minesweeper thinking goes.
// Everything here is yours to write; FirstCovered is only a placeholder
// so the server runs end to end.
package solver

import (
	"Schrodingers_Mine/protocol"
	"context"
)

type FirstCovered struct{}

func (FirstCovered) NewGame(rows, cols, mines int) {}

func (FirstCovered) NextMove(ctx context.Context, st *protocol.State) (int, int, error) {
	for r, row := range st.Board {
		for c, v := range row {
			if v == protocol.Covered {
				return r, c, nil
			}
		}
	}
	return 0, 0, ErrNoMove
}
