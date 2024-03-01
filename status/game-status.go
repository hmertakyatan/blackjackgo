package status

type GameStatus uint8

func (g GameStatus) String() string {
	switch g {
	case GameStatusConnected:
		return "CONNECTED"
	case GameStatusPlayerReady:
		return "PLAYER READY"
	case GameStatusDealing:
		return "DEALING"
	case GamestatusDone:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}

const (
	GameStatusConnected GameStatus = iota
	GameStatusPlayerReady
	GameStatusDealing
	GamestatusDone
)
