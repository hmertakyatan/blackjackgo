package structures

import "sync"

type Table struct {
	Lock     sync.RWMutex
	Seats    map[int]*Player
	MaxSeats int
}
