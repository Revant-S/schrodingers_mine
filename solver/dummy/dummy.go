package dummy

import (
	"Schrodingers_Mine/protocol"
	"Schrodingers_Mine/solver"
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
	return 0, 0, solver.ErrNoMove
}
