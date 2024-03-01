package main

import (
	"fmt"

	"github.com/hmertakyatan/blackjackgo/services"
)

func main() {
	deck := services.CreateDeck(2)
	fmt.Println("Deck:", deck)
	shuffledDeck := services.ShuffleDeck(deck)
	fmt.Println("Shuffled Deck: ", shuffledDeck)

}
