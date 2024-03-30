package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hmertakyatan/blackjackgo/handlers"
	"github.com/hmertakyatan/blackjackgo/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	playerService := services.NewPlayerService()
	jwtService := services.NewJwtService()
	server := services.NewServer()
	roomService := services.NewRoomService(server)
	jwthandler := handlers.NewJwtHandler(jwtService)
	roomhandler := handlers.NewRoomHandler(roomService)
	websockethandler := handlers.NewWebSocketHandler(jwtService, playerService, server)
	go server.RunServerLoop()

	r := gin.Default()
	r.GET("/ws/:roomID", websockethandler.HandleWebSocket)
	r.POST("/generatetoken", jwthandler.GenerateTokenFromPayload)
	r.POST("/createroom", roomhandler.HandleCreateRoom)
	r.GET("/roomlist", roomhandler.HandleGetRoomList)

	r.Run(":8080")
}
