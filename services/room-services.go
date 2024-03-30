package services

import (
	"github.com/google/uuid"
	"github.com/hmertakyatan/blackjackgo/dto"
)

type Room struct {
	ID        string             `json:"room_id"`
	Name      string             `json:"room_name"`
	Players   map[string]*Player `json:"players"`
	SeatNum   int                `json:"seat_number"`
	Available bool               `json:"available"`
	GameState *GameState
}

type RoomService struct {
	server *Server
}

func NewRoomService(s *Server) *RoomService {
	return &RoomService{
		server: s,
	}
}

func (rs *RoomService) CreateRoom(roomreq *dto.RoomRequest) *Room {
	roomID := uuid.New().String()
	table := NewTable(roomreq.SeatsNum)

	rs.server.Rooms[roomID] = &Room{
		ID:        roomID,
		Name:      roomreq.RoomName,
		SeatNum:   roomreq.SeatsNum,
		Players:   make(map[string]*Player),
		Available: true,
		GameState: NewGame(table),
	}
	return rs.server.Rooms[roomID]

}

func (rs *RoomService) GetAllRooms() []*dto.RoomRespose {
	var rooms []*dto.RoomRespose
	for _, r := range rs.server.Rooms {
		room := &dto.RoomRespose{
			RoomID:   r.ID,
			RoomName: r.Name,
			SeatsNum: r.SeatNum,
		}
		rooms = append(rooms, room)
	}
	return rooms
}
