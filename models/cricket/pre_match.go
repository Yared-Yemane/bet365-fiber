package cricket_models

type PrematchResponse struct {
    Success int `json:"success"`
    Results []struct {
        ID         string `json:"id"`
        SportID    string `json:"sport_id"`
        Time       int64  `json:"time"`
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
        Markets []Market `json:"markets"`
    } `json:"results"`
}


type Market struct {
    Name     string `json:"name"`
    Header   string `json:"header"`
    Odds     string `json:"odds"`
    Handicap string `json:"handicap"`
}

