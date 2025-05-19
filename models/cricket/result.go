package cricket_models

// type CricketResultResponse struct {
// 	Success int `json:"success"`
// 	Results []struct {
// 		ID         string `json:"id"`
// 		SportID    string `json:"sport_id"`
// 		Time       string `json:"time"`
// 		TimeStatus string `json:"time_status"` // "3" means match completed
// 		League     struct {
// 			ID   string `json:"id"`
// 			Name string `json:"name"`
// 			CC   string `json:"cc"`
// 		} `json:"league"`
// 		Home struct {
// 			ID      string `json:"id"`
// 			Name    string `json:"name"`
// 			ImageID string `json:"image_id"`
// 			CC      string `json:"cc"`
// 		} `json:"home"`
// 		Away struct {
// 			ID      string `json:"id"`
// 			Name    string `json:"name"`
// 			ImageID string `json:"image_id"`
// 			CC      string `json:"cc"`
// 		} `json:"away"`
// 		SS    string `json:"ss"` // Score format "117-217" (home-away)
// 		Extra struct {
// 			StadiumData struct {
// 				ID           string `json:"id"`
// 				Name         string `json:"name"`
// 				City         string `json:"city"`
// 				Country      string `json:"country"`
// 				Capacity     string `json:"capacity"`
// 				GoogleCoords string `json:"googlecoords"`
// 			} `json:"stadium_data"`
// 		} `json:"extra"`
// 		Bet365ID string `json:"bet365_id"`
// 	} `json:"results"`
// }

type ResultResponse struct {
	Success int `json:"success"`
	Results []struct {
		ID         string `json:"id"`
		SportID    string `json:"sport_id"`
		Time       string `json:"time"`
		TimeStatus string `json:"time_status"`
		League     struct {
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
		SS    string `json:"ss"`
		Extra struct {
			StadiumData struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				City         string `json:"city"`
				Country      string `json:"country"`
				Capacity     string `json:"capacity"`
				GoogleCoords string `json:"googlecoords"`
			} `json:"stadium_data"`
		} `json:"extra"`
	} `json:"results"`
}
