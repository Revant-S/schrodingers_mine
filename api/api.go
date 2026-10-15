package api

import (
	"Schrodingers_Mine/protocol"
	"Schrodingers_Mine/solver"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type Config struct {
	Name       string               // sent in READY
	MoveBudget time.Duration        // thinking time per move; keep below the board's move timeout
	NewSolver  func() solver.Solver // one solver per connection

	InsecureSkipVerify bool // accept any Origin; turn off if ever in prod
}

const (
	handshakeTimeout = 10 * time.Second
	writeTimeout     = 2 * time.Second
	readLimit        = 64 << 20
)

// Handler returns the http.HandlerFunc for the board's WebSocket.
func Handler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		})
		if err != nil {
			log.Printf("accept: %v", err)
			return
		}
		defer func(conn *websocket.Conn) {
			err := conn.CloseNow()
			if err != nil {
				return
			}
		}(conn)
		conn.SetReadLimit(readLimit)

		log.Printf("board connected from %s", r.RemoteAddr)
		err = run(r.Context(), conn, cfg)

		switch status := websocket.CloseStatus(err); {
		case err == nil:
			err := conn.Close(websocket.StatusNormalClosure, "")
			if err != nil {
				return
			}
			log.Printf("run finished")
		case status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway:
			log.Printf("board closed the connection")
		default:
			log.Printf("session ended: %v", err)
		}
	}
}

// run is one board session: handshake, then games until SUMMARY.
func run(ctx context.Context, conn *websocket.Conn, cfg Config) error {
	if err := handshake(ctx, conn, cfg.Name); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	s := cfg.NewSolver()
	currentGame := 0

	for {
		data, err := read(ctx, conn, 0)
		if err != nil {
			return err
		}

		var env protocol.Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			log.Printf("bad frame from board: %v", err)
			continue
		}

		switch env.Type {
		case protocol.TypeState:
			var st protocol.State
			if err := json.Unmarshal(data, &st); err != nil {
				return fmt.Errorf("decode STATE: %w", err)
			}

			if st.Result != protocol.Ongoing {
				logResult(&st)
				continue // never answer a finished game
			}
			if st.Game != currentGame {
				currentGame = st.Game
				s.NewGame(st.Rows, st.Cols, st.Mines)
			}

			moveCtx, cancel := context.WithTimeout(ctx, cfg.MoveBudget)
			row, col, err := s.NextMove(moveCtx, &st)
			cancel()
			if err != nil {
				return fmt.Errorf("game %d: solver: %w", st.Game, err)
			}

			mv := protocol.Move{Type: protocol.TypeMove, Game: st.Game, Row: row, Col: col}
			if err := write(ctx, conn, mv); err != nil {
				return err
			}

		case protocol.TypeError:
			var e protocol.Error
			_ = json.Unmarshal(data, &e)
			log.Printf("game %d: board rejected move: %s (%s)", e.Game, e.Code, e.Msg)
			// the board re-sends the unchanged STATE next, so nothing to do here

		case protocol.TypeSummary:
			var sum protocol.Summary
			if err := json.Unmarshal(data, &sum); err != nil {
				return fmt.Errorf("decode SUMMARY: %w", err)
			}
			rate := 0.0
			if sum.Played > 0 {
				rate = 100 * float64(sum.Won) / float64(sum.Played)
			}
			log.Printf("SUMMARY: played %d, won %d, lost %d, win rate %.2f%%, losses %v",
				sum.Played, sum.Won, sum.Lost, rate, sum.LossReasons)
			return nil

		default:
			log.Printf("ignoring unknown message type %q", env.Type)
		}
	}
}

func handshake(ctx context.Context, conn *websocket.Conn, name string) error {
	data, err := read(ctx, conn, handshakeTimeout)
	if err != nil {
		return err
	}
	var h protocol.Hello
	if err := json.Unmarshal(data, &h); err != nil || h.Type != protocol.TypeHello {
		return fmt.Errorf("expected HELLO, got %s", data)
	}
	if h.Protocol != protocol.Version {
		return fmt.Errorf("board speaks protocol %d, bot speaks %d", h.Protocol, protocol.Version)
	}
	return write(ctx, conn, protocol.Ready{Type: protocol.TypeReady, Name: name})
}

// read reads one text message. timeout 0 means wait as long as ctx allows.
func read(ctx context.Context, conn *websocket.Conn, timeout time.Duration) ([]byte, error) {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	typ, data, err := conn.Read(ctx)
	if err != nil {
		return nil, err
	}
	if typ != websocket.MessageText {
		return nil, errors.New("board sent a binary frame")
	}
	return data, nil
}

func write(ctx context.Context, conn *websocket.Conn, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return conn.Write(ctx, websocket.MessageText, data)
}

func logResult(st *protocol.State) {
	switch st.Result {
	case protocol.Won:
		log.Printf("game %d (seed %d): WON in %d moves", st.Game, st.Seed, st.Moves)
	case protocol.Lost:
		reason := "?"
		if st.Reason != nil {
			reason = *st.Reason
		}
		log.Printf("game %d (seed %d): LOST (%s) after %d moves, exploded at %v",
			st.Game, st.Seed, reason, st.Moves, st.Exploded)
	}
}
