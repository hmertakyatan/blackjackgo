package services

import (
	"math/rand"
	"time"
)

func CreateDeck(deckQuantity int) []Card {
	types := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	values := []string{"A", "K", "Q", "J", "10", "9", "8", "7", "6", "5", "4", "3", "2"}
	var deck []Card

	for i := 0; i < deckQuantity; i++ {
		for _, suit := range types {
			for _, value := range values {
				card := Card{Value: value, Type: suit}
				deck = append(deck, card)
			}
		}
	}
	return deck
}

func ShuffleDeck(deck []Card) []Card {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func DealCard(deck []Card) (Card, []Card) {
	card := deck[0]       //this card is top card of the deck
	return card, deck[1:] // return top card and REST of the deck
}

func DealCardAndUpdatePlayer(deck []Card, player *Player) []Card {
	card, restOfDeck := DealCard(deck)
	player.Hand = append(player.Hand, card)
	return restOfDeck
}

func DealCardToDealer(deck []Card, dealer *Dealer) []Card {
	card, restOfDeck := DealCard(deck)
	dealer.Hand = append(dealer.Hand, card)
	return restOfDeck
}
