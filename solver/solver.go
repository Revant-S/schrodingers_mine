package solver

import (
	"Schrodingers_Mine/protocol"
	"context"
	"errors"
)

type Solver interface {
	NewGame(rows, cols, mines int)
	NextMove(ctx context.Context, st *protocol.State) (row, col int, err error)
}

var ErrNoMove = errors.New("no covered cell")
