package poker

// CONSTANTS
// Figures = Rank integer part
const (
	Rank_zero float64 = iota
	Rank_op   float64 = iota
	Rank_tp   float64 = iota
	Rank_toak float64 = iota
	Rank_s    float64 = iota
	Rank_f    float64 = iota
	Rank_fh   float64 = iota
	Rank_foak float64 = iota
	Rank_sf   float64 = iota
)

type Hand struct {
	input string
	cards []Card
	rank  float64
}

type Card struct {
	value float64
	suit  string
}

type KindOccurence struct {
	value float64
	count int
}
