package poker

func IsStraightFlush(cards *[]Card) bool {
	return IsStraight(cards) && IsFlush(cards)
}

func IsFourOfAKind(kinds KindsOccurence) bool {
	for _, v := range kinds {
		if v == 4 {
			return true
		}
	}
	return false
}

func IsFullHouse(cards *[]Card) (bool, float64) {
	a, b, i := 1, 1, 1
	var d float64
	var highest, lowest int

	for ; i < len(*cards); i++ {
		if (*cards)[i].value == (*cards)[0].value {
			a++
		} else {
			i++
			break
		}
	}

	if a == 3 {
		highest = a
		d = float64((*cards)[0].value) / 10.0
		lowest = a
	}

	for ; i < len(*cards); i++ {
		if (*cards)[i].value == (*cards)[4].value {
			b++
		}
	}

	if b == 3 {
		highest = b
		lowest = a
		d = float64((*cards)[4].value) / 10.0
	}

	if highest == 3 && lowest == 2 {
		return true, d
	}

	return false, d
}

func IsFlush(cards *[]Card) bool {
	for i := 1; i < len((*cards))-1; i++ {
		if (*cards)[i].suit != (*cards)[0].suit {
			return false
		}
	}
	return true
}

func IsStraight(cards *[]Card) bool {
	hasAce := (*cards)[0].value == 14
	straight := false

	for i := 0; i < len(*cards)-1; i++ {
		if (*cards)[i+1].value == (*cards)[i].value-1 {
			straight = true
		} else {
			straight = false
			break
		}
	}

	if hasAce && !straight {
		(*cards)[0].value = 1
		straight = false
		for i := 0; i < len((*cards))-1; i++ {
			if (*cards)[i+1].value == (*cards)[i].value-1 {
				straight = true
			}
		}

		if !straight {
			(*cards)[0].value = 14
		}
	}

	return straight
}

func IsThreeOfAKind(kinds KindsOccurence) bool {
	for _, v := range kinds {
		if v == 3 {
			return true
		}
	}
	return false
}

func IsTwoPair(kinds KindsOccurence) bool {
	var countPairs int
	for _, v := range kinds {
		if v == 2 {
			countPairs++
		}
	}
	return countPairs == 2
}

func IsOnePair(kinds KindsOccurence) bool {
	var countPairs int
	for _, v := range kinds {
		if v == 2 {
			countPairs++
		}
	}
	return countPairs == 1
}
