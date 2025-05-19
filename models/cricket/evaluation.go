package cricket_models

type BetEvaluationRequest struct {
	Market    string `json:"market"`
	Selection string `json:"selection"`
	Handicap  string `json:"handicap,omitempty"`
	ScoreLine string `json:"score_line,omitempty"`
}
