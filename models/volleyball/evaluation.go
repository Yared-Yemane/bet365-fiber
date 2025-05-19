package volleyball_models

// BetEvaluationRequest represents the payload for evaluating a bet
// @Description Request body for evaluating a betting selection
type BetEvaluationRequest struct {
	Market    string `json:"market"`
	Selection string `json:"selection"`
	Handicap  string `json:"handicap,omitempty"`
	ScoreLine string `json:"score_line,omitempty"` // Add this for correct score
}

// BetSelection represents a concrete betting selection
// @Description Concrete betting selection with all required parameters
type BetSelection struct {
	Market    string  `json:"market"`
	Selection string  `json:"selection"`
	Odds      float64 `json:"odds"`
	Handicap  string  `json:"handicap,omitempty"`
	ScoreLine string  `json:"score_line,omitempty"`
}


