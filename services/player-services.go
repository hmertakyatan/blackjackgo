package services

import (
	"fmt"

	"github.com/hmertakyatan/blackjackgo/status"
	"github.com/hmertakyatan/blackjackgo/structures"
)

var players []*structures.Player

func CreatePlayerFromToken(tokenClaims map[string]interface{}) (*structures.Player, error) {

	username, ok := tokenClaims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("user id not found from token claims")
	}

	userID, ok := tokenClaims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user id not found from token claims")
	}
	balance, ok := tokenClaims["balance"].(float64)
	if !ok {
		return nil, fmt.Errorf("balance not found from token claims")
	}

	player := &structures.Player{
		Username:      username,
		ID:            userID,
		Balance:       balance,
		Bet:           0,
		CurrentAction: status.PlayerActionNone,
		GameStatus:    status.GameStatusConnected,
		TablePosition: -1,
	}

	players = append(players, player) //Pushed new created user to this array. We will fetch players from this array.
	return player, nil
}

func GetPlayerById(id string) (*structures.Player, error) {
	for _, player := range players {
		if player.ID == id {
			return player, nil
		}
	}
	return nil, fmt.Errorf("Player not found with ID: %s", id)
}
