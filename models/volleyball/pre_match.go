package volleyball_models

type PrematchResponse struct {
	Success int `json:"success"`
	Results []struct {
		FI      string `json:"FI"`
		EventID string `json:"event_id"`
		Main    struct {
			UpdatedAt string `json:"updated_at"`
			Key       string `json:"key"`
			Sp        struct {
				GameLines struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Odds []Odd  `json:"odds"`
				} `json:"game_lines"`
				CorrectSetScore struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Odds []Odd  `json:"odds"`
				} `json:"correct_set_score"`
				// Other markets can be added here
			} `json:"sp"`
		} `json:"main"`
		Others []struct {
			UpdatedAt string `json:"updated_at"`
			Sp        struct {
				Set1Lines struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Odds []Odd  `json:"odds"`
				} `json:"set_1_lines"`
				// Other sub-markets can be added here
			} `json:"sp"`
		} `json:"others"`
		Schedule struct {
			UpdatedAt string `json:"updated_at"`
			Key       string `json:"key"`
			Sp        struct {
				Main []struct {
					ID       string `json:"id"`
					Odds     string `json:"odds"`
					Name     string `json:"name"`
					Handicap string `json:"handicap"`
				} `json:"main"`
			} `json:"sp"`
		} `json:"schedule"`
	} `json:"results"`
}

type Odd struct {
	ID       string `json:"id"`
	Odds     string `json:"odds"`
	Name     string `json:"name"`
	Header   string `json:"header"`
	Handicap string `json:"handicap"`
}
