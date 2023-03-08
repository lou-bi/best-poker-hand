package poker

func IsStraightFlush(cards []Card, kinds *[]KindOccurence) float64 {
	// Check for IsStraight AFTER as it will mutate "kinds" in case of a "low ace straight"
	if IsFlush(cards) > 0 && IsStraight(kinds) > 0 {
		return Rank_sf
	}
	return Rank_zero
}

func IsFourOfAKind(kinds []KindOccurence) float64 {
	for _, kind := range kinds {
		if kind.count == 4 {
			return Rank_foak
		}
	}

	return Rank_zero
}

func IsFullHouse(kinds []KindOccurence) float64 {
	if IsThreeOfAKind(kinds) != 0 && IsOnePair(kinds) != 0 {
		return Rank_fh
	}
	return Rank_zero
}

func IsFlush(cards []Card) float64 {
	for i := 1; i < len((cards))-1; i++ {
		if cards[i].suit != cards[0].suit {
			return Rank_zero
		}
	}
	return Rank_f
}

func IsStraight(kindsOccurence *[]KindOccurence) float64 {
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

func IsThreeOfAKind(kinds []KindOccurence) float64 {
	for _, kind := range kinds {
		if kind.count == 3 {
			return Rank_toak
		}
	}
	return Rank_zero
}

func IsTwoPair(kinds []KindOccurence) float64 {
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

func IsOnePair(kinds []KindOccurence) float64 {
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
