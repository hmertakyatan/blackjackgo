package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/hmertakyatan/blackjackgo/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	jwtsecret := os.Getenv("JWTSecretKey")
	tokenExpireStr := os.Getenv("Access_token_expire")
	tokenExpire, err := time.ParseDuration(tokenExpireStr)
	if err != nil {
		fmt.Println("Error converting token expiration duration:", err)
		return
	}

	table := services.CreateTable(6)

	playerlist := services.NewPlayersList()
	//Players info
	playerPayloads := []map[string]interface{}{
		{
			"username": "player1",
			"user_id":  "9edfa439-c805-4a40-8d0d-24adfc3c47e6",
			"balance":  100.0,
		},
		{
			"username": "player2",
			"user_id":  "737fa074-2376-469b-8d4b-673ed249831a",
			"balance":  250.0,
		},
		{
			"username": "player3",
			"user_id":  "6ef58271-630c-412a-aa46-53e8e64dbb95",
			"balance":  450.0,
		},
	}
	//Create player from given jwt payloads
	for _, payload := range playerPayloads {
		tokenstring, err := services.GenerateToken(time.Duration(tokenExpire), payload, jwtsecret)
		if err != nil {
			fmt.Println("got an error when generating token:", err)
			continue
		}
		subClaims, err := services.ValidateToken(tokenstring, jwtsecret)
		if err != nil {
			fmt.Println("got an error when generating token:", err)
			continue
		}
		MappedSubClaims, ok := subClaims.(map[string]interface{})
		if !ok {
			fmt.Println("unexpected type assertion error")
			continue
		}

		player, err := services.CreatePlayerFromToken(MappedSubClaims)
		if err != nil {
			fmt.Println("error creating player:", err)
			continue
		}

		playerlist.AddPlayerToPlayerList(player)
	}

	players := playerlist.GetPlayerList()

	fmt.Println("Players created:")
	for _, player := range players {
		fmt.Println("Username:", player.Username)
		fmt.Println("UserID:", player.ID)
		fmt.Println("Balance:", player.Balance)
		fmt.Println("Table Position:", player.TablePosition)
		fmt.Println("Current Action:", player.CurrentAction.PlayerActionCase())
		fmt.Println("Game Status:", player.GameStatus.GameStatusCase())
		fmt.Println("************************************************")

	}
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	table.AddPlayerOnPosition(players[0].ID, rand.Intn(7))
	table.AddPlayerOnPosition(players[1].ID, rand.Intn(7))
	table.AddPlayerOnPosition(players[2].ID, rand.Intn(7))

	playersatthetable := table.PlayersAtTheTable()
	fmt.Println("Players at the table:")
	for _, player := range playersatthetable {
		fmt.Println("Username:", player.Username)
		fmt.Println("UserID:", player.ID)
		fmt.Println("Balance:", player.Balance)
		fmt.Println("Table Position:", player.TablePosition)
		fmt.Println("Current Action:", player.CurrentAction.PlayerActionCase())
		fmt.Println("Game Status:", player.GameStatus.GameStatusCase())
		fmt.Println()
	}
}
