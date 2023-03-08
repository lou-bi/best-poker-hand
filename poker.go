package poker

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// The whole idea is to get, for each hands,
// a rank based on figures as its integer parts,
// and its sorted card values representing the kind, at its decimal part, such as "0.x1x2x3x4x5"
//
// This will make comparisition trivial for retrieving the winner(s)
func BestHand(hands []string) ([]string, error) {

	if err := checkHandsFormat(&hands); err != nil {
		return []string{}, err
	}

	if len(hands) <= 1 {
		return hands, nil
	}

	parsedHands := parseHands(hands)

	higherRank := parsedHands[0].rank
	winners := []string{parsedHands[0].input}

	for i := 1; i < len(parsedHands) && parsedHands[i].rank == higherRank; i++ {
		winners = append(winners, parsedHands[i].input)
	}

	return winners, nil
}

/**
 * We process a basic regex on every hand to see if we can continue.
 * If this pass, no more error can rise in the program.
 */
func checkHandsFormat(hands *[]string) error {
	const CardSuits = "[♡♤♢♧]"
	const CardKinds = `(?:[2-9]|10|[JQKA])`

	rs := fmt.Sprintf(`^(?:%[1]v%[2]v ){4}%[1]v%[2]v$`, CardKinds, CardSuits)

	for _, hand := range *hands {
		res, err := regexp.MatchString(rs, hand)
		if !res || err != nil {
			return errors.New(fmt.Sprintf("Hand %v is invalid", hand))
		}
	}

	return nil
}

// ParseHands returns sorted details about hands.
// This operation is heavy but makes futures calculations quite simple
func parseHands(hands []string) []Hand {
	var parsedHands []Hand

	for _, hand := range hands {
		handSplit := strings.Split(hand, " ")

		cards := parseCards(handSplit)
		rank := getRank(cards)

		parsedHand := Hand{hand, cards, rank}
		parsedHands = append(parsedHands, parsedHand)
	}

	sort.Slice(parsedHands, func(a, b int) bool {
		if parsedHands[a].rank == parsedHands[b].rank {
			return parsedHands[a].cards[0].value > parsedHands[b].cards[0].value
		}
		return parsedHands[a].rank > parsedHands[b].rank
	})

	return parsedHands
}

// ParseCards returns sorted details about hands.
// This operation is heavy but makes futures calculations quite simple
func parseCards(rawHandSplit []string) []Card {
	var cards []Card

	for _, h := range rawHandSplit {
		// suit index
		// I did keep the icon so len() is longer
		si := len(h) - 3

		value := h[0:si]
		suit := h[si:]
		intValue := getCardIntValue(value)

		cards = append(cards, Card{intValue, suit})
	}

	// higher first, as always
	sort.Slice(cards, func(a, b int) bool {
		return cards[a].value > cards[b].value
	})

	return cards
}

// GetCardIntValue convert card's kind to float64
// such as ("2" -> "A") = (2 -> 14)
func getCardIntValue(target string) float64 {
	orderedValues := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	for i, kind := range orderedValues {
		if target == kind {
			/*
					+ 2 to match the card real value
				 	(2 == 2; J == 11)
					(A == 14 OR A == 1 if used in a low straight)
			*/
			return float64(i) + 2.0
		}
	}
	// Will never reach because of previous verifications
	return -1
}

// GetRank return a (0 -> 14) floating rank based
// on figures matching + paired cards values
func getRank(cards []Card) float64 {

	kindsOccurence := getKindsOccurence(cards)
	var rankIntegerPart float64

	rankIntegerPart = IsStraightFlush(cards, &kindsOccurence)

	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsFourOfAKind(kindsOccurence)
	}
	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsFullHouse(kindsOccurence)
	}
	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsFlush(cards)
	}
	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsStraight(&kindsOccurence)
	}
	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsThreeOfAKind(kindsOccurence)
	}
	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsTwoPair(kindsOccurence)
	}
	if rankIntegerPart == Rank_zero {
		rankIntegerPart = IsOnePair(kindsOccurence)
	}

	rankDecimalPart := getRankDecimalPart(kindsOccurence)

	return rankIntegerPart + rankDecimalPart
}

// GetKindsOccurence returns a sorted struct slice,
// such as { value: <card kind>, <count: occurence in hand>}
func getKindsOccurence(cards []Card) []KindOccurence {
	var mapKindsOccurence = make(map[float64]int)
	var kindsOccurence []KindOccurence

	for _, c := range cards {
		if _, ok := mapKindsOccurence[c.value]; ok {
			mapKindsOccurence[c.value]++
		} else {
			mapKindsOccurence[c.value] = 1
		}
	}

	for value, count := range mapKindsOccurence {
		kindsOccurence = append(kindsOccurence, KindOccurence{value, count})
	}

	sort.Slice(kindsOccurence, func(a, b int) bool {
		return kindsOccurence[a].value > kindsOccurence[b].value
	})

	sort.Slice(kindsOccurence, func(a, b int) bool {
		return kindsOccurence[a].count > kindsOccurence[b].count
	})

	return kindsOccurence
}

//
// GetRankDecimalPart returns a (0.0000000000 -> 0.1414141414),
// where each pair of decimal stores sorted cards value.
//
// Example:
// [{ 14, 2 } { 7, 1 } { 4, 3 }] outputs 0.041407
//
/////////////////////////////////////////////////////////////////////////////
// 0.041407 =>   |       04         |      14         |        07
//               |  triplet's rank  |  pair's rank    |  remainder's rank                      
/////////////////////////////////////////////////////////////////////////////
func getRankDecimalPart(kindsOccurence []KindOccurence) float64 {
	var divider, sum float64 = 100, 100

	for _, kind := range kindsOccurence {
		sum += kind.value
		sum *= 100
		divider *= 100
	}

	return sum / divider - 1
}
