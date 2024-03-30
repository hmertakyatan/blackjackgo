package services

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/hmertakyatan/blackjackgo/status"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type BroadcastTo struct {
	To      []Player
	Payload any
}
type BaseMessage struct {
	Token  string `json:"token"`
	Roomid string `json:"room_id"`
}

type SeatMessage struct {
	BaseMessage
	Position int `json:"position"`
}
type ReadyMessage struct {
	BaseMessage
	Betvalue float64 `json:"bet_value"`
}

type MessagePlayerAction struct {
	CurrentGameStatus status.GameStatus
	Action            status.PlayerAction
	Value             int
}

func (p *Player) WriteMessage(msg *Message) error {

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.Conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (p *Player) ReadMessage(s *Server) {
	defer func() {
		s.UnregisterPlayer <- p
	}()
	var msg Message
	for {

		if err := p.Conn.ReadJSON(&msg); err != nil { //handle message here and send to handleMessage method
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("err %v", err)
			}
			break
		}
		go s.handleMessage(&msg) //select according to message type and apply the service

	}

}
