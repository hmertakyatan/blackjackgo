package services

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/hmertakyatan/blackjackgo/status"
)

type Player struct {
	ID            string `json:"user_id"`
	Username      string `json:"username"`
	Conn          *websocket.Conn
	Balance       float64             `json:"balance"`
	Bet           float64             `json:"bet"`
	Hand          []Card              `json:"player_hand"`
	RoomID        string              `json:"room_id"`
	CurrentAction status.PlayerAction `json:"current_player_action"`
	GameStatus    status.GameStatus   `json:"game_status"`
	TablePosition int                 `json:"player_table_posion"`
	Message       chan *Message
}

type PlayerService struct{}

func NewPlayerService() *PlayerService {
	return &PlayerService{}
}

func (ps *PlayerService) CreatePlayerFromTokenClaimsAndConnection(tokenClaims map[string]interface{}, conn *websocket.Conn, roomid string) (*Player, error) {
	subject, ok := tokenClaims["sub"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("sub not found in token claims")
	}

	username, ok := subject["username"].(string)
	if !ok {
		return nil, fmt.Errorf("username not found in token claims")
	}

	userID, ok := subject["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id not found in token claims")
	}

	balance, ok := subject["balance"].(float64)
	if !ok {
		return nil, fmt.Errorf("balance not found in token claims")
	}

	player := &Player{
		Username:      username,
		ID:            userID,
		Balance:       balance,
		Bet:           0,
		Conn:          conn,
		RoomID:        roomid,
		CurrentAction: status.PlayerActionNone,
		GameStatus:    status.GameStatusConnected,
		TablePosition: -1,
	}
	return player, nil
}
