package poker

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)



func getCardIntValue(target string) float64 {
	orderedValues := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	for i, e := range orderedValues {
		if target == e {
			// + 2 to match the card real value
			// (2 == 2; J == 11)
			// (A == 14 OR A == 1 if used in a low straight)
			return float64(i) + 2.0
		}
	}
	// Will never reach because of previous verifications
	return -1
}

/**
 * Main program.
 */
func BestHand(hands []string) ([]string, error) {
	if err := checkHandsFormat(&hands); err != nil {
		return []string{""}, err
	}

	parsedHands := parseHands(hands)

	higherRank := parsedHands[0].rank
	winners := []string{parsedHands[0].input}

	for _, h := range parsedHands[1:] {
		if h.rank == higherRank {

			for i := 0; i < 5; i++ {
				v := h.cards[i].value
				x := parsedHands[0].cards[i].value
				if v > x {
					return []string{h.input}, nil
				}
				if v < x {
					return winners, nil
				}

			}

			winners = append(winners, h.input)

		}
	}

	return winners, nil
}

/**
 * We process a basic regex on every hand to see if we can continue.
 * If this pass, no more error can rise in the program.
 */
func checkHandsFormat(hands *[]string) error {
	const CardSuits = "[♡♤♢♧]"
	const CardValues = `(?:[2-9]|10|[JQKA])`

	rs := fmt.Sprintf(`^(?:%[1]v%[2]v ){4}%[1]v%[2]v$`, CardValues, CardSuits)

	for _, hand := range *hands {
		res, err := regexp.MatchString(rs, hand)
		if !res || err != nil {
			return errors.New("Hand " + hand + " is invalid.")
		}
	}

	return nil
}

/**
 * Parse hands in a more practical format to work with.
 * Also sort hands by value from higher to lower
 */
func parseHands(hands []string) []Hand {
	var parsedHands []Hand

	for _, hand := range hands {
		var cards []Card
		handSplit := strings.Split(hand, " ")

		for _, h := range handSplit {
			// suit index
			// I did keep the icon so len() is longer
			si := len(h) - 3

			value := h[0:si]
			suit := h[si:]
			intValue := getCardIntValue(value)

			cards = append(cards, Card{intValue, suit})
		}

		sort.Slice(cards, func(a, b int) bool {
			return cards[a].value > cards[b].value
		})
		cardsPtr := &cards
		parsedHands = append(parsedHands, Hand{hand, *cardsPtr, getRank(cardsPtr)})
	}

	sort.Slice(parsedHands, func(a, b int) bool {
		if parsedHands[a].rank == parsedHands[b].rank {
			return parsedHands[a].cards[0].value > parsedHands[b].cards[0].value
		}
		return parsedHands[a].rank > parsedHands[b].rank
	})

	return parsedHands
}

func getRank(cards *[]Card) float64 {
	// Todo: extract this
	var kindsMap = make(map[float64]int)
	var kindsOccurence []KindOccurence

	for _, c := range *cards {
		if _, ok := kindsMap[c.value]; ok {
			kindsMap[c.value]++
		} else {
			kindsMap[c.value] = 1
		}
	}

	for value, count := range kindsMap {
		kindsOccurence = append(kindsOccurence, KindOccurence{value, count})
	}
	/////////
	
	if IsStraightFlush(cards) {
		return 8
	}
	if IsFourOfAKind(kindsOccurence) {
		return 7
	}
	if res, d := IsFullHouse(*cards); res {
		return 6 + d
	}
	if IsFlush(*cards) {
		return 5
	}
	if IsStraight(cards) {
		return 4
	}
	if IsThreeOfAKind(kindsOccurence) {
		return 3
	}
	if res, d := IsTwoPair(kindsOccurence); res {
		return 2 + d
	}
	if IsOnePair(kindsOccurence) {
		return 1
	}
	return 0
}
