package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/cardrank/cardrank"
)

const (
	GameType    = cardrank.StudHiLo
	PlayerCount = 2
)

func main() {
	count := 30000

	heroHand := []cardrank.Card{
		cardrank.New(cardrank.Ace, cardrank.Spade),
		cardrank.New(cardrank.Ace, cardrank.Heart),
		cardrank.New(cardrank.Eight, cardrank.Heart),
		cardrank.New(cardrank.Two, cardrank.Heart),
	}
	villainHand := []cardrank.Card{
		cardrank.New(cardrank.Ace, cardrank.Club),
		cardrank.New(cardrank.Two, cardrank.Spade),
	}

	winCount := 0

	for i := 0; i < count; i++ {
		hiOrder, loOrder := stud(heroHand, villainHand)
		if len(hiOrder) > 0 && len(loOrder) > 0 {
			if hiOrder[0] == loOrder[0] && hiOrder[0] == 0 {
				winCount += 1
			}
		} else {
			if hiOrder[0] == 0 {
				winCount += 1
			}
		}
	}

	result := float64(winCount) / float64(count)

	fmt.Printf("result: %f\n", result)
}

func stud(heroHand, villainHand []cardrank.Card) ([]int, []int) {
	pockets, board := deal(heroHand, villainHand)

	evs := GameType.EvalPockets(pockets, board)

	for i := 0; i < PlayerCount; i++ {
		hi, lo := evs[i].Desc(false), evs[i].Desc(true)
		fmt.Printf("  %d: %b %b %b %s\n", i+1, pockets[i], hi.Best, hi.Unused, hi)
		fmt.Printf("                   %b %b %s\n", lo.Best, lo.Unused, lo)
	}
	hiOrder, hiPivot := cardrank.Order(evs, false)
	loOrder, loPivot := cardrank.Order(evs, true)
	hi := evs[hiOrder[0]].Desc(false)
	if hiPivot == 1 {
		fmt.Printf("Result: %d wins with %s, %b\n", hiOrder[0]+1, hi, hi.Best)
	} else {
		var s []string
		for i := 0; i < hiPivot; i++ {
			s = append(s, strconv.Itoa(hiOrder[i]+1))
		}
		fmt.Printf("Result: %s push with %s\n", strings.Join(s, ", "), hi)
	}
	if loPivot == 0 {
		fmt.Printf("        None\n")
	} else if loPivot == 1 {
		lo := evs[loOrder[0]].Desc(true)
		fmt.Printf("        %d wins with %s %b\n", loOrder[0]+1, lo, lo.Best)
	} else {
		var s []string
		for j := 0; j < loPivot; j++ {
			s = append(s, strconv.Itoa(loOrder[j]+1))
		}
		lo := evs[loOrder[0]].Desc(true)
		fmt.Printf("        %s push with %s\n", strings.Join(s, ", "), lo)
	}

	return hiOrder, loOrder
}

func deal(heroHand, villainHand []cardrank.Card) ([][]cardrank.Card, []cardrank.Card) {
	seed := time.Now().UnixNano()
	// note: use a better pseudo-random number generator
	r := rand.New(rand.NewSource(seed))

	fmt.Printf("------ StudHiLo %d ------\n", seed)

	//pockets, board := GameType.Deal(r, 3, PlayerCount)
	//

	desc := GameType.Desc()
	deckType := desc.Deck
	deckType.Exclude()

	cards := deckType.Exclude(heroHand, villainHand)
	deck := cardrank.DeckOf(cards...)
	deck.Shuffle(r, 3)

	dealer := cardrank.NewDealer(desc, deck, PlayerCount)

	for dealer.Next() {
	}
	pockets := dealer.Runs[0].Pockets

	board := dealer.Runs[0].Hi

	hero := removeElements(pockets[0], 0, len(heroHand))
	heroHand = append(heroHand, hero...)

	villain := removeElements(pockets[1], 0, len(villainHand))
	villainHand = append(villainHand, villain...)

	return [][]cardrank.Card{heroHand, villainHand}, board
}

func removeElements(arr []cardrank.Card, start, count int) []cardrank.Card {
	end := start + count
	if end > len(arr) {
		end = len(arr)
	}

	return append(arr[:start], arr[end:]...)
}
