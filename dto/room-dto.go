package dto

type RoomRequest struct {
	RoomName string `json:"room_name"`
	SeatsNum int    `json:"seats_num"`
}

type RoomRespose struct {
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
	SeatsNum int    `json:"seats_num"`
}
