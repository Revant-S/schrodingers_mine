package protocol

import "encoding/json"

const Version = 1

const (
	TypeHello   = "HELLO"
	TypeReady   = "READY"
	TypeState   = "STATE"
	TypeMove    = "MOVE"
	TypeError   = "ERROR"
	TypeSummary = "SUMMARY"
)

const (
	Covered = -2 // still hidden
	Mine    = -3 // only shown after a loss to a mine
)

const (
	Ongoing = 0
	Won     = 1
	Lost    = 2
)

type Envelope struct {
	Type string `json:"type"`
}

type Hello struct {
	Type     string `json:"type"`
	Protocol int    `json:"protocol"`
}

type Ready struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type State struct {
	Type     string  `json:"type"`
	Game     int     `json:"game"`
	Games    int     `json:"games"`
	Seed     uint32  `json:"seed"`
	Rows     int     `json:"rows"`
	Cols     int     `json:"cols"`
	Mines    int     `json:"mines"`
	Board    [][]int `json:"board"` // Board[row][col]
	Result   int     `json:"result"`
	Reason   *string `json:"reason"` // nil unless Result == Lost
	Moves    int     `json:"moves"`
	Last     []int   `json:"last"`     // [row, col] or nil
	Exploded []int   `json:"exploded"` // [row, col] or nil
}

type Move struct {
	Type string `json:"type"`
	Game int    `json:"game"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}

type Error struct {
	Type string `json:"type"`
	Game int    `json:"game"`
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type Summary struct {
	Type        string          `json:"type"`
	Played      int             `json:"played"`
	Won         int             `json:"won"`
	Lost        int             `json:"lost"`
	LossReasons map[string]int  `json:"lossReasons"`
	TotalMoves  int             `json:"totalMoves"`
	ElapsedMs   int64           `json:"elapsedMs"`
	Games       json.RawMessage `json:"games"` // per-game rows; kept raw, can be huge
}
