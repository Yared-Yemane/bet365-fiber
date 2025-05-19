package volleyball_models

type SetScore struct {
	Home string `json:"home"`
	Away string `json:"away"`
}

type ResultResponse struct {
	Success int `json:"success"`
	Results []struct {
		ID         string `json:"id"`
		SportID    string `json:"sport_id"`
		Time       string `json:"time"`
		TimeStatus string `json:"time_status"`
		// Scores     map[string]SetScore `json:"scores"`
		League struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			CC   string `json:"cc"`
		} `json:"league"`
		Home struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			ImageID string `json:"image_id"`
			CC      string `json:"cc"`
		} `json:"home"`
		Away struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			ImageID string `json:"image_id"`
			CC      string `json:"cc"`
		} `json:"away"`
		SS     string              `json:"ss"`
		Scores map[string]SetScore `json:"scores"`
		Stats  struct {
			PointsWonOnServe []string `json:"points_won_on_serve"`
			LongestStreak    []string `json:"longest_streak"`
		} `json:"stats"`
		Events []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"events"`
		Extra struct {
			HomePos    string `json:"home_pos"`
			AwayPos    string `json:"away_pos"`
			BestOfSets string `json:"bestofsets"`
			Round      string `json:"round"`
		} `json:"extra"`
		InplayCreatedAt string `json:"inplay_created_at"`
		InplayUpdatedAt string `json:"inplay_updated_at"`
		ConfirmedAt     string `json:"confirmed_at"`
		Bet365ID        string `json:"bet365_id"`
	} `json:"results"`
}
