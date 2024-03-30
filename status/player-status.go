package status

type PlayerAction uint8

func (pa PlayerAction) PlayerActionCase() string {
	switch pa {
	case PlayerActionNone:
		return "NONE"
	case PlayerActionHit:
		return "HIT"
	case PlayerActionStand:
		return "STAND"
	case PlayerActionDouble:
		return "DOUBLE"
	case PlayerActionSplit:
		return "SPLIT"
	case PlayerActionBet:
		return "BETTING"
	case PlayerActionReady:
		return "READY"
	case PlayerActionDealing:
		return "DEALING"
	case PlayerActionSitOnSeat:
		return "SEAT"
	default:
		return "INVALID"
	}
}

const (
	PlayerActionNone PlayerAction = iota
	PlayerActionHit
	PlayerActionStand
	PlayerActionDouble
	PlayerActionSplit
	PlayerActionBet
	PlayerActionReady
	PlayerActionDealing
	PlayerActionSitOnSeat
)
