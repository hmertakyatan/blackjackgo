package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hmertakyatan/blackjackgo/services"
)

type WebSocketHandler struct {
	upgrader      websocket.Upgrader
	jwtService    *services.JwtService
	playerService *services.PlayerService
	server        *services.Server
}

func NewWebSocketHandler(jwtService *services.JwtService, playerService *services.PlayerService, server *services.Server) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		jwtService:    jwtService,
		playerService: playerService,
		server:        server,
	}
}

func (wh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	roomid := c.Param("roomID")
	if err := wh.jwtService.ValidateToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Token not valid.": err.Error()})
		return
	}
	claims, err := wh.jwtService.ExtractAllClaimsFromToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
		return
	}
	conn, err := wh.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player, err := wh.playerService.CreatePlayerFromTokenClaimsAndConnection(claims, conn, roomid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Player cannot create from token. ERROR: ": err.Error()})
		return
	}
	defer conn.Close()
	fmt.Println("New connection: Player ID:", player.ID, conn.RemoteAddr().String())
	wh.server.RegisterPlayer <- player

	go player.WriteMessage()

	player.ReadMessage(wh.server)

}
