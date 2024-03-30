package services

type Game struct {
	PlayerHand []Card
	DealerHand []Card
	Deck       []Card
	GameOver   bool
	PlayerWins bool
}

func (g *Game) ScoreCalculator(hand []Card) (int, bool) {
	score := 0
	aces := 0     //aces quantity
	soft := false //if player has aces and his score below than 21, it is soft hand, like 2/12 , 5/15 etc.

	for _, card := range hand {
		switch {
		case card.Value == "A":
			aces++
			score += 11
		case card.Value == "K", card.Value == "Q", card.Value == "J":
			score += 10
		case card.Value == "10":
			score += 10
		default:
			score += int(card.Value[0] - '0')
		}
	}

	for score > 21 && aces > 0 {
		score -= 10
		aces--
	}

	if aces > 0 && score <= 21 {
		soft = true
	}

	return score, soft
}
