package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hmertakyatan/blackjackgo/status"
)

type Server struct {
	Players          map[string]*Player
	Rooms            map[string]*Room
	RegisterPlayer   chan *Player
	UnregisterPlayer chan *Player
	Broadcast        chan *Message
	MsgChan          chan *Message
	jwtService       *JwtService
}

func NewServer() *Server {

	s := &Server{
		Players:          make(map[string]*Player),
		Rooms:            make(map[string]*Room),
		RegisterPlayer:   make(chan *Player, 10),
		UnregisterPlayer: make(chan *Player),
		Broadcast:        make(chan *Message, 100),
		MsgChan:          make(chan *Message, 100),
	}
	return s

}

func (s *Server) RunServerLoop() {
	for {
		select {
		case p := <-s.RegisterPlayer:
			if _, ok := s.Rooms[p.RoomID]; ok {
				r := s.Rooms[p.RoomID]

				if _, ok := r.Players[p.ID]; !ok {
					r.Players[p.ID] = p
				}

			}

		case p := <-s.UnregisterPlayer:
			if _, ok := s.Rooms[p.RoomID]; ok {
				r := s.Rooms[p.ID]

				if _, ok := r.Players[p.ID]; ok {
					r.Players[p.ID] = p
					delete(s.Rooms[p.ID].Players, p.ID)
				}

			}

		case msg := <-s.MsgChan:
			go func() {
				if err := s.handleMessage(msg); err != nil {
					fmt.Println("handle message error.")
				}
			}()
		}
	}
}

func (s *Server) handleMessage(msg *Message) error {
	switch msg.Type {
	case status.PlayerActionSitOnSeat.PlayerActionCase():
		var seatMsg SeatMessage
		if err := json.Unmarshal(msg.Data, &seatMsg); err != nil {

			return err
		}
		return s.handleSeatMessage(seatMsg)
	case status.PlayerActionReady.PlayerActionCase():
		var readyMsg ReadyMessage
		if err := json.Unmarshal(msg.Data, &readyMsg); err != nil {
			customerr := errors.New("Unexpected datas.")
			return customerr
		}
		return s.handleReadyMessage(readyMsg)

	default:
		return errors.New("Unexpected message.")

	}
}

func (s *Server) handleSeatMessage(seatMsg SeatMessage) error {
	room, ok := s.Rooms[seatMsg.Roomid]
	if !ok {
		return errors.New("Room not found")
	}
	player, ok := room.Players[seatMsg.Roomid]
	if !ok {
		return errors.New("Player not found")
	}
	return room.GameState.table.AddPlayerOnPosition(player.ID, seatMsg.Position)
}

func (s *Server) handleReadyMessage(readyMsg ReadyMessage) error {
	room, ok := s.Rooms[readyMsg.Roomid]
	if !ok {
		return errors.New("Room not found")
	}
	playerid, err := s.jwtService.ExtractPlayerIdClaimFromToken(readyMsg.Token)
	if err != nil {
		return err
	}
	player, ok := room.Players[playerid]
	if !ok {
		return errors.New("Player not found")
	}
	if err := room.GameState.table.CollectBetFromPlayerOnTheTable(player.ID, readyMsg.Betvalue); err != nil { // set player ready in this method
		return err
	}
	room.GameState.SetStatus(status.GameStatusPlayerReady)

	readyPlayersCount := room.GameState.countReadyPlayers()

	if readyPlayersCount == 1 {
		room.GameState.InitiateShuffleAndDeal(1)
		return nil
	}

	time.AfterFunc(15*time.Second, func() {
		if room.GameState.countReadyPlayers() == readyPlayersCount {
			room.GameState.InitiateShuffleAndDeal(1)
		}
	})

	return nil

}
