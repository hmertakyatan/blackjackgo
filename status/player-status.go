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
)
