package poker

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const CardSuits = "[♡♤♢♧]"
const CardValues = "[2-9JQKA]"
const OrderedValues = "AKQJ98765432"

/**
 * Main program.
 */
func BestHand(hands []string) (string, error) {
	if err := checkHandsFormat(&hands); err != nil {
		return "", err
	}

	parsedHands := parseHands(&hands)

	for _, parsedHand := range parsedHands {
		isFourOfAKind(&parsedHand)
	}

	fmt.Printf("%v", parsedHands)
	return "", nil
}

/**
 * We process a basic regex on every hand to see if we can continue.
 * If this pass, no more error can rise in the program.
 */
func checkHandsFormat(hands *[]string) error {
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
func parseHands(hands *[]string) []Hand {
	var parsedHands []Hand

	for _, hand := range *hands {
		handSplit := strings.Split(hand, " ")

		sort.Slice(handSplit, func(a, b int) bool {
			aValue := []rune(handSplit[a])[0]
			bValue := []rune(handSplit[b])[0]
			return strings.IndexRune(OrderedValues, aValue) < strings.IndexRune(OrderedValues, bValue)
		})

		parsedHands = append(parsedHands, Hand{
			hand:          hand,
			formattedHand: handSplit,
		})
	}

	return parsedHands
}

// ///////////////////////////////////////////////////////////////////////////
// At this point, it is obvious that working plain string will make code unreadable.
// We should introduce some kind of structure to store a "rank" with a hand.
// We will probably follow the ranking system from the wikipedia page,
// i.e an int from 1 to 10, 1 being "high card" and 10 being "five of a kind".
//
// As "joker" is not a card from the test file, we will omit the 10th rank,
// going 1 to 9.
// ///////////////////////////////////////////////////////////////////////////
type Hand struct {
	hand          string
	formattedHand []string // this might be deleted when I figure out how
	rank          uint8
}

// ///////////////////////////////////////////////////////////////////////////
// We then write functions to detect all possible figures.
//
// We will use regex to detect figures based on repetitions,
// i.e flush and any pairs.
// (Yes, I love regexp, if it wasn't clear)
//
// They are all structured the same:
// - First we detect what we search for in a capturing group (suit OR value)
// - We then search for the captured char anywhere else, n - 1 times.
// (n-1 because we already matched it once in the captured group)
// ///////////////////////////////////////////////////////////////////////////

func isStraightFlush(h *Hand) bool {
	return isStraight(h) && isFlush(h)
}

func isFourOfAKind(h *Hand) bool {
	rs := fmt.Sprintf(`(%v)(?:.+\1){3}`, CardValues)
	r := regexp.MustCompile(rs)
	return r.MatchString(h.hand)
}

func isFullHouse(h *Hand) bool {
	return isThreeOfAKind(h) && isOnePair(h)
}

func isFlush(h *Hand) bool {
	rs := fmt.Sprintf(`(%v)(?:.+\1){4}`, CardSuits)
	r := regexp.MustCompile(rs)
	return r.MatchString(h.hand)
}

func isStraight(h *Hand) bool {
	return false
}

func isThreeOfAKind(h *Hand) bool {
	rs := fmt.Sprintf(`(%v)(?:.+\1){2}`, CardValues)
	r := regexp.MustCompile(rs)
	return r.MatchString(h.hand)
}

/**
 * This regex will match two pairs even if they are the same.
 * But because we will search for "isFourOfAKind" before,
 * we will never go there with 2 same pairs.
 */
func isTwoPair(h *Hand) bool {
	rs := fmt.Sprintf(`(?:.*?(%v).*?\1){2}`, CardValues)
	r := regexp.MustCompile(rs)
	return r.MatchString(h.hand)
}

/**
 * Same comment as above
 */
func isOnePair(h *Hand) bool {
	rs := fmt.Sprintf(`(%v).+\1`, CardValues)
	r := regexp.MustCompile(rs)
	return r.MatchString(h.hand)
}
