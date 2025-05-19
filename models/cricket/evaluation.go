package cricket_models

// type EvaluationResult struct {
//     Selection    BetSelection `json:"selection"`
//     ActualResult string       `json:"actual_result"`
//     Outcome      string       `json:"outcome"` // "won", "lost", "push", "void"
//     Description  string       `json:"description"`
// }

type BetEvaluationRequest struct {
	Market    string `json:"market"`
	Selection string `json:"selection"`
	Handicap  string `json:"handicap,omitempty"`
	ScoreLine string `json:"score_line,omitempty"`
}
