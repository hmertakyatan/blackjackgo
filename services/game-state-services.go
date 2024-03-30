package services

import (
	"fmt"

	"github.com/hmertakyatan/blackjackgo/status"
)

type GameState struct {
	player              *Player
	currentStatus       *AtomicInt
	currentPlayerAction *AtomicInt
	currentDealer       *AtomicInt
	currentPlayerTurn   *AtomicInt
	table               *TableService
	PlayersList         *PlayersList
}

func NewGame(table *TableService) *GameState {
	g := &GameState{
		currentStatus:       NewAtomicInt(int32(status.GameStatusConnected)),
		PlayersList:         NewPlayersList(),
		currentPlayerAction: NewAtomicInt(0),
		currentDealer:       NewAtomicInt(0),
		currentPlayerTurn:   NewAtomicInt(0),
		table:               table,
	}
	return g
}
func (g *GameState) PlayersAtTheGame() []*Player {
	players := g.PlayersList.GetPlayerList()
	return players
}

func (g *GameState) DealerAtTheGame() *Dealer {
	dealer := g.table.Dealer
	return dealer
}

func (g *GameState) CanTakeAction(player *Player) bool {
	if status.GameStatus(g.currentStatus.GetAtomicInt()) != status.GameStatusInitialDeal {
		return true
	} else {
		return false
	}
}

func (g *GameState) IncrementToNextPlayer() {
	player, err := g.table.GetPlayerAfter(g.player.ID)
	if err != nil {
		panic(err)
	}

	if g.PlayersList.Len()-1 == int(g.currentPlayerTurn.GetAtomicInt()) {
		g.currentPlayerTurn.SetAtomicInt(0)
		return
	}
	g.currentPlayerTurn.IncrementAtomicIntValue()

	fmt.Println("the next player on the table is:", player.TablePosition)
}

func (g *GameState) AdvanceToNexRound() {
	g.currentPlayerAction.SetAtomicInt(int32(status.PlayerActionNone))

	if status.GameStatus(g.currentStatus.GetAtomicInt()) == status.GamestatusDone {
		g.SetReady()
		return
	}
}

func (g *GameState) SetReady() {
	tablePosition := g.PlayersList.GetPlayerIndexFromPlayerList(g.player.ID)
	g.table.AddPlayerOnPosition(g.player.ID, tablePosition)
	g.SetStatus(status.GameStatusPlayerReady)
}

func (g *GameState) TakeAction(action status.PlayerAction, value int) (MessagePlayerAction, error) {

	g.currentPlayerAction.SetAtomicInt((int32)(action))

	g.IncrementToNextPlayer()

	a := MessagePlayerAction{
		Action:            action,
		CurrentGameStatus: status.GameStatus(g.currentStatus.GetAtomicInt()),
		Value:             value,
	}

	return a, nil
}

func (g *GameState) isCanGameStart() bool {
	playersAtTheTable := g.table.PlayersAtTheTable()

	for _, player := range playersAtTheTable {
		if status.GameStatus(player.CurrentAction) != status.GameStatus(status.PlayerActionReady) {
			return false
		}
	}
	return true
}

func (g *GameState) AppendReadyPlayersToPlayersAtGameList(readyplayers []*Player) {
	g.PlayersList.Lock.Lock()
	defer g.PlayersList.Lock.Unlock()

	g.PlayersList.List = append(g.PlayersList.List, readyplayers...) //wtf i need learn this sht
}

func (g *GameState) InitiateShuffleAndDeal(deckQuantity int) {
	if !g.isCanGameStart() {
		fmt.Println("Game cannot start because players are not ready.")
		return
	}

	g.AppendReadyPlayersToPlayersAtGameList(g.table.PlayersAtTheTable())

	deck := CreateDeck(deckQuantity)

	deck = ShuffleDeck(deck)

	for _, player := range g.PlayersAtTheGame() {
		deck = DealCardAndUpdatePlayer(deck, player)
	}

	deck = DealCardToDealer(deck, g.table.Dealer)

	for _, player := range g.PlayersAtTheGame() {
		deck = DealCardAndUpdatePlayer(deck, player)
	}

	deck = DealCardToDealer(deck, g.table.Dealer)

	g.SetStatus(status.GameStatusInitialDeal)
}

func (g *GameState) SetStatus(s status.GameStatus) {

	if status.GameStatus(g.currentStatus.GetAtomicInt()) != s {
		g.currentStatus.SetAtomicInt(int32(s))
	}
}

func (g *GameState) isAllPlayersReady() bool {
	for _, player := range g.table.PlayersAtTheTable() {
		if player.CurrentAction != status.PlayerActionReady {
			return false
		}
	}
	return true
}

func (g *GameState) countReadyPlayers() int {
	count := 0
	for _, player := range g.table.PlayersAtTheTable() {
		if player.CurrentAction == status.PlayerActionReady {
			count++
		}
	}
	return count
}
