package structures

import (
	"github.com/hmertakyatan/blackjackgo/status"
)

type Player struct {
	ID            string              `json:"user_id"`
	Username      string              `json:"username"`
	Balance       float64             `json:"balance"`
	Bet           float64             `json:"bet"`
	CurrentAction status.PlayerAction `json:"current_player_action"`
	GameStatus    status.GameStatus   `json:"game_status"`
	TablePosition int                 `json:"player_table_posion"`
}
