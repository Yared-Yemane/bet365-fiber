package volleyball_models

type SetScore struct {
	Home string `json:"home"`
	Away string `json:"away"`
}

type ResultResponse struct {
	Success int `json:"success"`
	Results []struct {
		// ... other fields ...
		SS     string `json:"ss"`
		Scores map[string]struct {
			Home string `json:"home"`
			Away string `json:"away"`
		} `json:"scores"`
		// ... other fields ...
	} `json:"results"`
}
