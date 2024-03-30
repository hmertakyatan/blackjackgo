package api

import "github.com/hmertakyatan/blackjackgo/services"

type ApiServer struct {
	Port string
	Game *services.GameState
}

func NewAPIServer(port string, game *services.GameState) *ApiServer {
	return &ApiServer{
		Game: game,
		Port: port,
	}
}
