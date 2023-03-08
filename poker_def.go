package poker

type Hands = []Hand

type Hand struct {
	input string
	cards []Card
	rank  float64
}

type Card struct {
	value float64
	suit  string
}

// {value: count}
type KindsOccurence = map[float64]int
