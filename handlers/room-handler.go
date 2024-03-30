package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hmertakyatan/blackjackgo/dto"
	"github.com/hmertakyatan/blackjackgo/services"
)

type RoomHandler struct {
	roomService *services.RoomService
}

func NewRoomHandler(roomService *services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

func (rh *RoomHandler) HandleCreateRoom(c *gin.Context) {
	var roomrequest dto.RoomRequest
	if err := c.ShouldBindJSON(&roomrequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdRoom := rh.roomService.CreateRoom(&roomrequest)

	c.JSON(http.StatusOK, createdRoom)
}

func (rh *RoomHandler) HandleGetRoomList(c *gin.Context) {
	roomList := rh.roomService.GetAllRooms()
	c.JSON(http.StatusOK, roomList)
}
