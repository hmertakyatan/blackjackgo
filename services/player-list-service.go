package services

import (
	"fmt"
	"sort"
	"sync"
)

type PlayersList struct {
	Lock sync.RWMutex
	List []*Player
}

func NewPlayersList() *PlayersList {
	return &PlayersList{
		List: []*Player{},
	}
}

func (p *PlayersList) GetPlayerList() []*Player {
	p.Lock.RLock()

	defer p.Lock.RUnlock()

	return p.List
}

func (p *PlayersList) GetPlayerFromPlayerListByIndex(index any) *Player {
	p.Lock.RLock()

	defer p.Lock.RUnlock()

	var i int
	switch v := index.(type) {
	case int:
		i = v
	case int32:
		i = int(v)
	}

	if len(p.List)-1 < i {
		panic("The given index is too high")
	}

	return p.List[i]
}

// *************
func (p *PlayersList) Len() int {
	return len(p.List)
}

func (p *PlayersList) Less(i, j int) bool {
	return p.List[i].ID < p.List[j].ID
}

func (p *PlayersList) Swap(i, j int) {
	p.List[i], p.List[j] = p.List[j], p.List[i]
}

// *************
func (p *PlayersList) AddPlayerToPlayerList(player *Player) *Player {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	p.List = append(p.List, player)
	sort.Sort(p)
	return player
}

func (p *PlayersList) AddDealerToPlayerListAsPlayer(dealer *Player) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	p.List = append(p.List, dealer)
	sort.Sort(p)
}

//*************

func (p *PlayersList) GetPlayerIndexFromPlayerList(userid string) int {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	for i := 0; i < len(p.List); i++ {
		if userid == p.List[i].ID {
			return i
		}
	}
	return -1
}

func (p *PlayersList) GetPlayerById(id string) (*Player, error) {
	p.Lock.RLock()
	defer p.Lock.RUnlock()
	for _, player := range p.List {
		if player.ID == id {
			return player, nil
		}
	}
	return nil, fmt.Errorf("Player not found with ID: %s", id)
}
