package poker

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

/**
 * Main program.
 */
func BestHand(hands []string) (string, error) {
	if err := checkHandsFormat(&hands); err != nil {
		return "", err
	}
	parsedHands := parseHands(&hands)
	fmt.Printf("%v", parsedHands)
	return "", nil
}

/**
 * We process a basic regex on every hand to see if we can continue.
 * If this pass, no more error can rise in the program.
 */
func checkHandsFormat(hands *[]string) error {
	validHandRegex := `^(?:[2-9JQKA][♡♤♢♧] ){4}[2-9JQKA][♡♤♢♧]$`

	for _, hand := range *hands {
		res, err := regexp.MatchString(validHandRegex, hand)
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
func parseHands(hands *[]string) [][]string {
	var parsedHands [][]string

	orders := "AKQJ98765432"

	for _, hand := range *hands {
		handSplit := strings.Split(hand, " ")

		sort.Slice(handSplit, func(a, b int) bool {
			aValue := []rune(handSplit[a])[0]
			bValue := []rune(handSplit[b])[0]
			return strings.IndexRune(orders, aValue) < strings.IndexRune(orders, bValue)
		})

		parsedHands = append(parsedHands, handSplit)
	}

	return parsedHands
}
