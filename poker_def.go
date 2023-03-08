package poker

type Hands = []Hand

type Hand struct {
	input string
	cards []Card
	rank  float64
}

type Card struct {
	value int //
	suit  string
}

type KindsOccurence = map[int]int
