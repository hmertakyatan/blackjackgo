package services

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/hmertakyatan/blackjackgo/status"
)

type TableService struct {
	Lock       sync.RWMutex
	Seats      map[int]interface{}
	MaxSeats   int
	PlayerList *PlayersList
	Dealer     *Dealer
}
type Dealer struct {
	Player
	ID            string
	Username      string
	Balance       float64 `json:"balance"`
	Bet           float64 `json:"bet"`
	Hand          []Card
	CurrentAction status.PlayerAction `json:"current_player_action"`
	GameStatus    status.GameStatus   `json:"game_status"`
	TablePosition int                 `json:"player_table_posion"`
}
type Card struct {
	Type  string
	Value string
}

func CreateNewDealer() *Dealer {
	dealer := &Dealer{
		ID:            uuid.New().String(),
		Username:      "Dealer",
		Hand:          []Card{},
		CurrentAction: status.PlayerActionReady,
		Balance:       0,
		TablePosition: 0,
	}
	return dealer
}

func NewTable(maxSeats int) *TableService {
	table := &TableService{
		Seats:      make(map[int]interface{}),
		MaxSeats:   maxSeats,
		PlayerList: NewPlayersList(),
		Dealer:     CreateNewDealer(),
	}
	table.AddDealerOnPosition(table.Dealer)
	return table
}

func (t *TableService) PlayersAtTheTable() []*Player {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	//we check every single seat and selecting players. Because we have a dealer on the seat.
	players := []*Player{}
	for _, seat := range t.Seats {
		if player, ok := seat.(*Player); ok {
			players = append(players, player)
		}
	}

	return players
}

func (t *TableService) RemovePlayerFromTableById(userid string) error {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	for i := 0; i < t.MaxSeats; i++ {
		player, ok := t.Seats[i].(*Player)
		if ok {
			if player.ID == userid {
				delete(t.Seats, i)
				return nil
			}
		}
	}

	return fmt.Errorf("player (%s) not on the table", userid)
}

func (t *TableService) GetPlayerFromTableById(userid string) (*Player, error) {
	for i := 0; i < t.MaxSeats; i++ {
		player, ok := t.Seats[i].(*Player)
		if ok {
			if player.ID == userid {
				return player, nil
			}
		}
	}
	return nil, fmt.Errorf("Player not found with ID: %s", userid)
}

func (t *TableService) SetPlayerStatus(userid string, newStatus status.GameStatus) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	p, err := t.GetPlayerFromTableById(userid)
	if err != nil {
		panic(err)
	}
	p.GameStatus = newStatus
}
func (t *TableService) isPositionOccupied(position int) bool {
	_, occupied := t.Seats[position]
	if occupied {
		return true
	} else {
		return false
	}
}
func (t *TableService) AddPlayerOnPosition(userid string, position int) error {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if len(t.Seats) == t.MaxSeats {
		return fmt.Errorf("player table is full")
	}

	if t.isPositionOccupied(position) {
		return fmt.Errorf("position %d is already occupied", position)
	}

	player, err := t.PlayerList.GetPlayerById(userid)
	if err != nil {
		return err
	}

	t.Seats[position] = player
	player.TablePosition = position
	player.CurrentAction = status.PlayerActionBet

	return nil
}

func (t *TableService) AddDealerOnPosition(dealer *Dealer) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	dealer.CurrentAction = status.PlayerActionReady
	t.Seats[0] = dealer
	dealer.TablePosition = 0
	dealer.CurrentAction = status.PlayerActionReady
	dealer.GameStatus = status.GameStatusPlayerReady

}

func (t *TableService) GetPlayerAfter(userid string) (*Player, error) {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	currentPlayer, err := t.GetPlayerFromTableById(userid)
	if err != nil {
		return nil, err
	}

	i := currentPlayer.TablePosition + 1
	for {
		nextPlayer, ok := t.Seats[i].(*Player)
		if ok {
			return nextPlayer, nil
		}

		i++
		if i >= t.MaxSeats {
			i = 0
		}

		// Tüm oyuncuların dolaşıldığını kontrol et
		if i == currentPlayer.TablePosition {
			return nil, fmt.Errorf("%s is the only player on the table", userid)
		}
	}
}

func (t *TableService) CollectBetFromPlayerOnTheTable(userid string, value float64) error {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	if value <= 0 {
		return fmt.Errorf("bet cannot be zero or negative ")
	}
	player, err := t.PlayerList.GetPlayerById(userid)
	if err != nil {
		return fmt.Errorf("Player cannot find from table player list")
	}

	if player.Balance < value {
		return fmt.Errorf("not enough balance player: %s", player.Username)
	}
	player.Bet = value
	player.Balance -= value
	player.CurrentAction = status.PlayerActionReady
	player.GameStatus = status.GameStatusPlayerReady

	return nil
}

func (t *TableService) JoinToTableList(player *Player) {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	t.PlayerList.List = append(t.PlayerList.List, player)

}
