package poker

/*
	This file make sense only with the comprehension of
	structures []Card and []KindOccurence.
	The fact that they are all sorted make algorythm much more concise.
*/

func isStraightFlush(cards []Card, kinds *[]KindOccurence) float64 {
	// Check for IsStraight AFTER as it will mutate "kinds" in case of a "low ace straight"
	if isFlush(cards) > 0 && isStraight(kinds) > 0 {
		return Rank_sf
	}
	return Rank_zero
}

func isFourOfAKind(kinds []KindOccurence) float64 {
	for _, kind := range kinds {
		if kind.count == 4 {
			return Rank_foak
		}
	}

	return Rank_zero
}

func isFullHouse(kinds []KindOccurence) float64 {
	if isThreeOfAKind(kinds) != 0 && isOnePair(kinds) != 0 {
		return Rank_fh
	}
	return Rank_zero
}

func isFlush(cards []Card) float64 {
	for i := 1; i < len((cards))-1; i++ {
		if cards[i].suit != cards[0].suit {
			return Rank_zero
		}
	}
	return Rank_f
}

//
// This one is a bit huge mostly because we need to check for "low start ace straight" figure,
// with the need of mutating ace's value.
// This way, we can geet the proper decimal rank for this figure.
//
func isStraight(kindsOccurence *[]KindOccurence) float64 {
	if len(*kindsOccurence) != 5 {
		return Rank_zero
	}

	hasAce := (*kindsOccurence)[0].value == 14
	straight := false

	for i := 0; i < len(*kindsOccurence)-1; i++ {
		if (*kindsOccurence)[i].value == (*kindsOccurence)[i+1].value+1 {
			straight = true
		} else {
			straight = false
			break
		}
	}

	if hasAce && !straight {
		straight = false

		// We move the low ace to the end of the slice,
		// by slicing its first element
		// and appending a {value: 1 count: 1} to KindsOccurence
		// We will save it if we have a straight
		newKindsOccurence := (*kindsOccurence)[1:]
		newKindsOccurence = append(newKindsOccurence, KindOccurence{1, 1})

		for i := 0; i < len(newKindsOccurence)-1; i++ {

			if (newKindsOccurence)[i].value == (newKindsOccurence)[i+1].value+1 {
				straight = true
			} else {
				return Rank_zero
			}

		}

		if straight {
			// We have a low ace straight, we need to save the new value of the ace
			// in order to get decimal rank right
			*kindsOccurence = newKindsOccurence
			return Rank_s
		}
	}

	if straight {
		return Rank_s
	}
	return Rank_zero
}

func isThreeOfAKind(kinds []KindOccurence) float64 {
	for _, kind := range kinds {
		if kind.count == 3 {
			return Rank_toak
		}
	}
	return Rank_zero
}

func isTwoPair(kinds []KindOccurence) float64 {
	var nbPair int

	for _, kind := range kinds {
		if kind.count == 2 {
			nbPair++
		}
	}

	if nbPair == 2 {
		return Rank_tp
	}
	return Rank_zero
}

func isOnePair(kinds []KindOccurence) float64 {
	var countPairs int
	for _, kind := range kinds {
		if kind.count == 2 {
			countPairs++
		}
	}
	if countPairs == 1 {
		return Rank_op
	}
	return Rank_zero
}
