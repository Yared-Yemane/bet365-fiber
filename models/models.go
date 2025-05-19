package models

type AvailableSelection struct {
	Market     string `json:"market"`
	Selections []struct {
		Name     string `json:"name"`
		Odds     string `json:"odds"`
		Handicap string `json:"handicap,omitempty"`
	} `json:"selections"`
}

// EvaluationResult represents the outcome of a bet evaluation
// @Description Result of evaluating a betting selection
type EvaluationResult struct {
	Selection    BetSelection `json:"selection"`
	ActualResult string       `json:"actual_result"`
	Outcome      string       `json:"outcome"`
	Description  string       `json:"description"`
}

type BetSelection struct {
	Market    string `json:"market"`
	Selection string `json:"selection"`
	Odds      string `json:"odds"`
	Handicap  string `json:"handicap,omitempty"`
	ScoreLine string `json:"score_line,omitempty"`
}
