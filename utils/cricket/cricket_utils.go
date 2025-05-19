package cricket_utils

import (
	"bet365-fiber-sim/models"
	cricket_models "bet365-fiber-sim/models/cricket"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var PrematchData cricket_models.PrematchResponse
var ResultData cricket_models.ResultResponse

func InitCricketHandlers() {
	var err error
	PrematchData, err = ReadCricketPrematchData("data/cricket_prematch.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load prematch data: %v", err))
	}

	ResultData, err = ReadCricketResultData("data/cricket_result.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load result data: %v", err))
	}
}

// ReadPrematchData reads and parses prematch JSON data

func ReadCricketPrematchData(filename string) (cricket_models.PrematchResponse, error) {
	var data cricket_models.PrematchResponse
	file, err := os.Open(filename)
	if err != nil {
		return data, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return data, fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return data, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return data, nil
}

// ReadResultData reads and parses result JSON data
func ReadCricketResultData(filename string) (cricket_models.ResultResponse, error) {
	var data cricket_models.ResultResponse
	file, err := os.Open(filename)
	if err != nil {
		return data, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return data, fmt.Errorf("failed to read file: %v", err)
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return data, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return data, nil
}

// CreateSelectionFromPrematch creates a bet selection from prematch data
func CreateCricketSelectionFromPrematch(data cricket_models.PrematchResponse, market, selection string, handicap ...string) models.BetSelection {
	for _, result := range data.Results {
		for _, m := range result.Markets {
			if m.Name == market && (m.Header == selection || m.Name == selection) && (len(handicap) == 0 || m.Handicap == handicap[0]) {
				odds := m.Odds
				return models.BetSelection{
					Market:    market,
					Selection: selection,
					Odds:      odds,
					Handicap:  m.Handicap,
				}
			}
		}
	}
	return models.BetSelection{}
}

// EvaluateSelection evaluates a bet selection against the result data
func EvaluateCricketSelection(selection models.BetSelection, resultData cricket_models.ResultResponse) models.EvaluationResult {
	if len(resultData.Results) == 0 {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "no result data",
			Outcome:      "void",
			Description:  "No result data available",
		}
	}

	result := resultData.Results[0]
	scores := strings.Split(result.SS, "-")
	if len(scores) != 2 {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "invalid score format",
			Outcome:      "void",
			Description:  "Score format is invalid",
		}
	}

	homeRuns, err1 := strconv.Atoi(scores[0])
	awayRuns, err2 := strconv.Atoi(scores[1])
	if err1 != nil || err2 != nil {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "invalid score values",
			Outcome:      "void",
			Description:  "Could not parse score values",
		}
	}

	switch selection.Market {
	case "Match Winner":
		return EvaluateCricketMatchWinner(selection, homeRuns, awayRuns)
	case "Total Runs":
		return EvaluateCricketTotalRuns(selection, homeRuns, awayRuns)
	case "Correct Score":
		return EvaluateCricketCorrectScore(selection, homeRuns, awayRuns)
	case "Double Chance":
		return EvaluateCricketDoubleChance(selection, homeRuns, awayRuns)
	default:
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "unknown market",
			Outcome:      "void",
			Description:  "Unknown market type",
		}
	}
}

// EvaluateMatchWinner evaluates a Match Winner bet
func EvaluateCricketMatchWinner(selection models.BetSelection, homeRuns, awayRuns int) models.EvaluationResult {
	var winner string
	if homeRuns > awayRuns {
		winner = "1"
	} else if awayRuns > homeRuns {
		winner = "2"
	} else {
		winner = "X" // Tie
	}

	outcome := "lost"
	if selection.Selection == winner {
		outcome = "won"
	}

	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: fmt.Sprintf("%d-%d (%s)", homeRuns, awayRuns, winner),
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s, actual winner was %s", selection.Selection, winner),
	}
}

// EvaluateTotalRuns evaluates a Total Runs bet
func EvaluateCricketTotalRuns(selection models.BetSelection, homeRuns, awayRuns int) models.EvaluationResult {
	parts := strings.Fields(selection.Handicap)
	if len(parts) < 2 {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "invalid total format",
			Outcome:      "void",
			Description:  fmt.Sprintf("Invalid total format: '%s'", selection.Handicap),
		}
	}

	target, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "invalid total value",
			Outcome:      "void",
			Description:  fmt.Sprintf("Failed to parse total value '%s': %v", parts[1], err),
		}
	}

	outcome := "lost"
	total := float64(homeRuns + awayRuns)

	switch parts[0] {
	case "O":
		if total > target {
			outcome = "won"
		} else if total == target {
			outcome = "push"
		}
	case "U":
		if total < target {
			outcome = "won"
		} else if total == target {
			outcome = "push"
		}
	default:
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "invalid total type",
			Outcome:      "void",
			Description:  fmt.Sprintf("Invalid total type '%s' (must be O/U)", parts[0]),
		}
	}

	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: fmt.Sprintf("%d runs (Home: %d, Away: %d)", homeRuns+awayRuns, homeRuns, awayRuns),
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s %s (target %.1f), actual was %.0f", parts[0], selection.Handicap, target, total),
	}
}

// EvaluateCorrectScore evaluates a Correct Score bet
func EvaluateCricketCorrectScore(selection models.BetSelection, homeRuns, awayRuns int) models.EvaluationResult {
	actualScore := fmt.Sprintf("%d-%d", homeRuns, awayRuns)
	outcome := "lost"
	if selection.ScoreLine == actualScore {
		outcome = "won"
	}

	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: actualScore,
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s, actual was %s", selection.ScoreLine, actualScore),
	}
}

// EvaluateDoubleChance evaluates a Double Chance bet
func EvaluateCricketDoubleChance(selection models.BetSelection, homeRuns, awayRuns int) models.EvaluationResult {
	var actualOutcome string
	if homeRuns > awayRuns {
		actualOutcome = "1"
	} else if homeRuns == awayRuns {
		actualOutcome = "X"
	} else {
		actualOutcome = "2"
	}

	outcome := "lost"
	switch selection.Selection {
	case "1X":
		if actualOutcome == "1" || actualOutcome == "X" {
			outcome = "won"
		}
	case "12":
		if actualOutcome == "1" || actualOutcome == "2" {
			outcome = "won"
		}
	case "X2":
		if actualOutcome == "X" || actualOutcome == "2" {
			outcome = "won"
		}
	}

	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: fmt.Sprintf("%d-%d (%s)", homeRuns, awayRuns, actualOutcome),
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s, actual was %s", selection.Selection, actualOutcome),
	}
}

// Get1X2Selections returns available Match Winner selections
func GetCricket1X2Selections() models.AvailableSelection {
	selections := []struct {
		Name     string `json:"name"`
		Odds     string `json:"odds"`
		Handicap string `json:"handicap,omitempty"`
	}{}

	for _, result := range PrematchData.Results {
		for _, market := range result.Markets {
			if market.Name == "Match Winner" {
				odds := market.Odds
				selections = append(selections, struct {
					Name     string `json:"name"`
					Odds     string `json:"odds"`
					Handicap string `json:"handicap,omitempty"`
				}{
					Name:     market.Header,
					Odds:     odds,
					Handicap: market.Handicap,
				})
			}
		}
	}

	return models.AvailableSelection{
		Market:     "Match Winner (1X2)",
		Selections: selections,
	}
}

// GetTotalRunsSelections returns available Total Runs selections
func GetCricketTotalRunsSelections() models.AvailableSelection {
	selections := []struct {
		Name     string `json:"name"`
		Odds     string `json:"odds"`
		Handicap string `json:"handicap,omitempty"`
	}{}

	for _, result := range PrematchData.Results {
		for _, market := range result.Markets {
			if market.Name == "Total Runs" {
				odds := market.Odds
				selections = append(selections, struct {
					Name     string `json:"name"`
					Odds     string `json:"odds"`
					Handicap string `json:"handicap,omitempty"`
				}{
					Name:     market.Header,
					Odds:     odds,
					Handicap: market.Handicap,
				})
			}
		}
	}

	return models.AvailableSelection{
		Market:     "Total Runs",
		Selections: selections,
	}
}

// GetDoubleChanceSelections returns available Double Chance selections
func GetCricketDoubleChanceSelections() models.AvailableSelection {
	return models.AvailableSelection{
		Market: "Double Chance",
		Selections: []struct {
			Name     string `json:"name"`
			Odds     string `json:"odds"`
			Handicap string `json:"handicap,omitempty"`
		}{
			{Name: "1X", Odds: "1.10"},
			{Name: "12", Odds: "1.05"},
			{Name: "X2", Odds: "1.20"},
		},
	}
}
