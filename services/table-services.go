package services

import (
	"fmt"
	"sync"

	"github.com/hmertakyatan/blackjackgo/status"
	"github.com/hmertakyatan/blackjackgo/structures"
)

type Table struct {
	Lock     sync.RWMutex
	Seats    map[int]*structures.Player
	MaxSeats int
}

func CreateTable(maxSeats int) *Table {
	return &Table{
		Seats:    make(map[int]*structures.Player),
		MaxSeats: maxSeats,
	}
}

func (t *Table) PlayersAtTheTable() []*structures.Player {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	players := []*structures.Player{}
	for i := 0; i < t.MaxSeats; i++ {
		player, ok := t.Seats[i]
		if ok {
			players = append(players, player)
		}
	}

	return players
}

func (t *Table) RemovePlayerFromTableById(userid string) error {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	for i := 0; i < t.MaxSeats; i++ {
		player, ok := t.Seats[i]
		if ok {
			if player.ID == userid {
				delete(t.Seats, i)
				return nil
			}
		}
	}

	return fmt.Errorf("player (%s) not on the table", userid)
}

func (t *Table) GetPlayerById(userid string) (*structures.Player, error) {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	for i := 0; i < t.MaxSeats; i++ {
		player, ok := t.Seats[i]
		if ok {
			if player.ID == userid {
				return player, nil
			}
		}

	}
	return nil, fmt.Errorf("Player not found with ID: %s", userid)
}

func (t *Table) SetPlayerStatus(userid string, newStatus status.GameStatus) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	p, err := t.GetPlayerById(userid)
	if err != nil {
		panic(err)
	}
	p.GameStatus = newStatus
}

func (t *Table) AddPlayerOnPosition(userid string, position int) error {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if len(t.Seats) == t.MaxSeats {
		return fmt.Errorf("player table is full")
	}

	player, err := GetPlayerById(userid)
	if err != nil {
		return err
	}
	player.TablePosition = position
	player.GameStatus = status.GameStatusPlayerReady

	t.Seats[position] = player

	return nil
}

func (t *Table) GetPlayerAfter(userid string) (*structures.Player, error) {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	currentPlayer, err := t.GetPlayerById(userid)
	if err != nil {
		return nil, err
	}

	i := currentPlayer.TablePosition + 1
	for {
		nextPlayer, ok := t.Seats[i]
		if nextPlayer == currentPlayer {
			return nil, fmt.Errorf("%s is the only player on the table", userid)
		}
		if ok {
			return nextPlayer, nil
		}

		i++
		if t.MaxSeats <= i {
			i = 0
		}
	}
}
