package services

import (
	"math/rand"
	"time"

	"github.com/hmertakyatan/blackjackgo/structures"
)

func CreateDeck(deckQuantity int) []structures.Card {
	types := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	values := []string{"A", "K", "Q", "J", "10", "9", "8", "7", "6", "5", "4", "3", "2"}
	var deck []structures.Card

	for i := 0; i < deckQuantity; i++ {
		for _, suit := range types {
			for _, value := range values {
				card := structures.Card{Value: value, Type: suit}
				deck = append(deck, card)
			}
		}
	}
	return deck
}

func ShuffleDeck(deck []structures.Card) []structures.Card {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func DealCard(deck []structures.Card) (structures.Card, []structures.Card) {
	card := deck[0]       //this card is top card of the deck
	return card, deck[1:] // return top card and REST of the deck
}
