package volleyball_utils

import (
	"bet365-fiber-sim/models"
	volleyball_models "bet365-fiber-sim/models/volleyball"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var PrematchData volleyball_models.PrematchResponse
var ResultData volleyball_models.ResultResponse

func InitHandlers(app *fiber.App) {
	// Load data at startup
	var err error
	PrematchData, err = ReadPrematchData("data/volleyball_prematch.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load prematch data: %v", err))
	}

	ResultData, err = ReadResultData("data/volleyball_result.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load result data: %v", err))
	}
}

// func ReadJSONFile[T any](path string, target *T) error {
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal(file, target)
// }

func ReadPrematchData(filename string) (volleyball_models.PrematchResponse, error) {
	var data volleyball_models.PrematchResponse

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

func ReadResultData(filename string) (volleyball_models.ResultResponse, error) {
	var data volleyball_models.ResultResponse

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

func CreateSelectionFromPrematch(data volleyball_models.PrematchResponse, market, header string, handicap ...string) models.BetSelection {
	for _, result := range data.Results {
		// Check main markets
		for _, odd := range result.Main.Sp.GameLines.Odds {
			if odd.Name == market && odd.Header == header && (len(handicap) == 0 || odd.Handicap == handicap[0]) {
				// odds, _ := strconv.ParseFloat(odd.Odds, 64)
				odds := fmt.Sprintf("%s", odd.Odds)

				return models.BetSelection{
					Market:    market,
					Selection: header,
					Odds:      odds,
					Handicap:  odd.Handicap,
				}
			}
		}

		// Check schedule
		for _, odd := range result.Schedule.Sp.Main {
			if odd.Name == market && (len(handicap) == 0 || odd.Handicap == handicap[0]) {
				odds := fmt.Sprintf("%s", odd.Odds)
				return models.BetSelection{
					Market:    market,
					Selection: header,
					Odds:      odds,
					Handicap:  odd.Handicap,
				}
			}
		}
	}
	return models.BetSelection{}
}

func CreateCorrectScoreSelection(data volleyball_models.PrematchResponse, header, score string) models.BetSelection {
	for _, result := range data.Results {
		for _, odd := range result.Main.Sp.CorrectSetScore.Odds {
			if odd.Header == header && odd.Name == score {
				odds := fmt.Sprintf("%s", odd.Odds)
				return models.BetSelection{
					Market:    "Correct Set Score",
					Selection: score,
					Odds:      odds,
				}
			}
		}
	}
	return models.BetSelection{}
}

func CreateDoubleChanceSelection(data volleyball_models.PrematchResponse, combo string) models.BetSelection {
	// In a real implementation, you would look up the actual odds for these combinations
	// This is a simplified version for demonstration
	return models.BetSelection{
		Market:    "Double Chance",
		Selection: combo,
		Odds:      "1.25", // Example odds
	}
}

func EvaluateSelection(selection models.BetSelection, resultData volleyball_models.ResultResponse) models.EvaluationResult {
	if len(resultData.Results) == 0 {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "no result data",
			Outcome:      "void",
			Description:  "No result data available",
		}
	}

	result := resultData.Results[0]
	totalPoints := CalculateTotalPoints(result.Scores)
	homeSets, awaySets := ParseSetScore(result.SS)

	switch selection.Market {
	case "Winner": // 1X2 Market
		return Evaluate1X2(selection, homeSets, awaySets)
	case "Total":
		return EvaluateTotal(selection, totalPoints)
	case "Correct Set Score":
		return EvaluateCorrectScore(selection, homeSets, awaySets)
	case "Double Chance":
		return EvaluateDoubleChance(selection, homeSets, awaySets)
	default:
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "unknown market",
			Outcome:      "void",
			Description:  "Unknown market type",
		}
	}
}

func Evaluate1X2(selection models.BetSelection, homeSets, awaySets int) models.EvaluationResult {
	var actualOutcome string
	if homeSets > awaySets {
		actualOutcome = "1"
	} else if homeSets == awaySets {
		actualOutcome = "X"
	} else {
		actualOutcome = "2"
	}

	outcome := "lost"
	if selection.Selection == actualOutcome {
		outcome = "won"
	}

	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: fmt.Sprintf("%d-%d (%s)", homeSets, awaySets, actualOutcome),
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s, actual result was %s", selection.Selection, actualOutcome),
	}
}

func EvaluateTotal(selection models.BetSelection, totalPoints int) models.EvaluationResult {
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
	actualTotal := float64(totalPoints)

	switch parts[0] {
	case "O":
		if actualTotal > target {
			outcome = "won"
		} else if actualTotal == target {
			outcome = "push"
		}
	case "U":
		if actualTotal < target {
			outcome = "won"
		} else if actualTotal == target {
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
		ActualResult: fmt.Sprintf("%d points", totalPoints),
		Outcome:      outcome,
		Description: fmt.Sprintf("Selected %s %s (target %.1f), actual was %d",
			parts[0], selection.Handicap, target, totalPoints),
	}
}

func EvaluateWinner(selection models.BetSelection, homeSets, awaySets int) models.EvaluationResult {
	var winner string
	if homeSets > awaySets {
		winner = "1"
	} else {
		winner = "2"
	}

	outcome := "lost"
	if selection.Selection == winner {
		outcome = "won"
	}

	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: fmt.Sprintf("%d-%d", homeSets, awaySets),
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s to win, actual result was %d-%d", selection.Selection, homeSets, awaySets),
	}
}

func EvaluateHandicap(selection models.BetSelection, homeSets, awaySets int, scores map[string]volleyball_models.SetScore) models.EvaluationResult {
	// First check if we even have a handicap value
	if selection.Handicap == "" {
		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: "no handicap provided",
			Outcome:      "void",
			Description:  "No handicap value was provided for this selection",
		}
	}

	// Try to parse the handicap value
	var handicapValue float64
	var err error

	// Check if it's a simple number (like "-1.5")
	if strings.HasPrefix(selection.Handicap, "-") || strings.HasPrefix(selection.Handicap, "+") {
		handicapValue, err = strconv.ParseFloat(selection.Handicap, 64)
		if err != nil {
			return models.EvaluationResult{
				Selection:    selection,
				ActualResult: "invalid handicap format",
				Outcome:      "void",
				Description:  fmt.Sprintf("Failed to parse handicap value '%s': %v", selection.Handicap, err),
			}
		}
	} else {
		// Handle cases like "O 177.5" or "U 177.5"
		parts := strings.Fields(selection.Handicap)
		if len(parts) < 2 {
			return models.EvaluationResult{
				Selection:    selection,
				ActualResult: "invalid handicap format",
				Outcome:      "void",
				Description:  fmt.Sprintf("Handicap format is invalid: '%s'", selection.Handicap),
			}
		}

		handicapValue, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return models.EvaluationResult{
				Selection:    selection,
				ActualResult: "invalid handicap value",
				Outcome:      "void",
				Description:  fmt.Sprintf("Failed to parse handicap value '%s': %v", parts[1], err),
			}
		}
	}

	// For set handicap (like -1.5)
	if strings.Contains(selection.Handicap, ".") {
		adjustedHome := float64(homeSets) - handicapValue
		outcome := "lost"
		if adjustedHome > float64(awaySets) {
			outcome = "won"
		}

		return models.EvaluationResult{
			Selection:    selection,
			ActualResult: fmt.Sprintf("%d-%d (handicap adjusted: %.1f-%d)", homeSets, awaySets, adjustedHome, awaySets),
			Outcome:      outcome,
			Description:  fmt.Sprintf("Selected %s with handicap %s", selection.Selection, selection.Handicap),
		}
	}

	// For point handicap (would need to calculate total points per team)
	return models.EvaluationResult{
		Selection:    selection,
		ActualResult: "point handicap not implemented",
		Outcome:      "void",
		Description:  "Point handicap evaluation not implemented",
	}
}

func EvaluateCorrectScore(selection models.BetSelection, homeSets, awaySets int) models.EvaluationResult {
	actualScore := fmt.Sprintf("%d-%d", homeSets, awaySets)
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

func EvaluateDoubleChance(selection models.BetSelection, homeSets, awaySets int) models.EvaluationResult {
	var actualOutcome string
	if homeSets > awaySets {
		actualOutcome = "1"
	} else if homeSets == awaySets {
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
		ActualResult: fmt.Sprintf("%d-%d (%s)", homeSets, awaySets, actualOutcome),
		Outcome:      outcome,
		Description:  fmt.Sprintf("Selected %s, actual was %s", selection.Selection, actualOutcome),
	}
}

func CalculateTotalPoints(scores map[string]volleyball_models.SetScore) int {
	total := 0
	for _, set := range scores {
		home, _ := strconv.Atoi(set.Home)
		away, _ := strconv.Atoi(set.Away)
		total += home + away
	}
	return total
}

func ParseSetScore(ss string) (int, int) {
	parts := strings.Split(ss, "-")
	if len(parts) != 2 {
		return 0, 0
	}
	home, _ := strconv.Atoi(parts[0])
	away, _ := strconv.Atoi(parts[1])
	return home, away
}

func CreateSelectionFromRequest(req volleyball_models.BetEvaluationRequest) models.BetSelection {
	switch req.Market {
	case "Winner", "Total":
		return FindSelectionInPrematch(req)
	case "Correct Set Score":
		return FindCorrectScoreSelection(req)
	case "Double Chance":
		return models.BetSelection{
			Market:    req.Market,
			Selection: req.Selection,
			Odds:      GetDoubleChanceOdds(req.Selection),
		}
	default:
		return models.BetSelection{}
	}
}

func FindSelectionInPrematch(req volleyball_models.BetEvaluationRequest) models.BetSelection {
	for _, result := range PrematchData.Results {
		// Check game lines
		for _, odd := range result.Main.Sp.GameLines.Odds {
			if odd.Name == req.Market &&
				(odd.Header == req.Selection || odd.Name == req.Selection) &&
				(req.Handicap == "" || odd.Handicap == req.Handicap) {
				odds := fmt.Sprintf("%s", odd.Odds)
				return models.BetSelection{
					Market:    req.Market,
					Selection: req.Selection,
					Odds:      odds,
					Handicap:  odd.Handicap,
				}
			}
		}

		// Check schedule
		for _, odd := range result.Schedule.Sp.Main {
			if odd.Name == req.Market &&
				(req.Handicap == "" || odd.Handicap == req.Handicap) {
				odds := fmt.Sprintf("%s", odd.Odds)
				return models.BetSelection{
					Market:    req.Market,
					Selection: req.Selection,
					Odds:      odds,
					Handicap:  odd.Handicap,
				}
			}
		}
	}
	return models.BetSelection{}
}

func Get1X2Selections() models.AvailableSelection {
	selections := []struct {
		Name     string  `json:"name"`
		Odds     string `json:"odds"`
		Handicap string  `json:"handicap,omitempty"`
	}{}

	for _, result := range PrematchData.Results {
		// 1. Check main GameLines
		// for _, odd := range result.Main.Sp.GameLines.Odds {
		// 	if (odd.Header == "1" || odd.Header == "2") && odd.Odds != "" && odd.Name == "Winner" {
		// 		if odds, err := strconv.ParseFloat(odd.Odds, 64); err == nil {
		// 			selections = append(selections, struct {
		// 				Name     string  `json:"name"`
		// 				Odds     float64 `json:"odds"`
		// 				Handicap string  `json:"handicap,omitempty"`
		// 			}{
		// 				Name:     odd.Header,
		// 				Odds:     odds,
		// 				Handicap: odd.Handicap,
		// 			})
		// 		}
		// 	}
		// }

		// 2. Check Schedule
		// for _, odd := range result.Schedule.Sp.Main {
		// 	if odd.Name == "Winner" && odd.Odds != "" {
		// 		if odds, err := strconv.ParseFloat(odd.Odds, 64); err == nil {
		// 			selectionName := "1"       // default to home
		// 			if odd.ID == "666717703" { // Away team ID
		// 				selectionName = "2"
		// 			}
		// 			selections = append(selections, struct {
		// 				Name     string  `json:"name"`
		// 				Odds     float64 `json:"odds"`
		// 				Handicap string  `json:"handicap,omitempty"`
		// 			}{
		// 				Name:     selectionName,
		// 				Odds:     odds,
		// 				Handicap: odd.Handicap,
		// 			})
		// 		}
		// 	}
		// }

		// 3. Check Others sections
		for _, other := range result.Others {
			// Check Set1Lines in Others
			if other.Sp.Set1Lines.Odds != nil {
				for _, odd := range other.Sp.Set1Lines.Odds {
					if odd.Name == "Winner" && odd.Odds != "" {
						odds := fmt.Sprintf("%s", odd.Odds)
						{
							selections = append(selections, struct {
								Name     string `json:"name"`
								Odds     string `json:"odds"`
								Handicap string `json:"handicap,omitempty"`
							}{
								Name:     odd.Header, // "1" or "2"
								Odds:     odds,
								Handicap: odd.Handicap,
							})
						}
					}
				}
			}

			// Add checks for other markets in Others if needed
			// e.g., other.Sp.Set1ToGoToExtraPoints, etc.
		}
	}

	// Remove duplicates if any
	uniqueSelections := removeDuplicateSelections(selections)

	return models.AvailableSelection{
		Market:     "Winner (1X2)",
		Selections: uniqueSelections,
	}
}

func removeDuplicateSelections(selections []struct {
	Name     string `json:"name"`
	Odds     string `json:"odds"`
	Handicap string `json:"handicap,omitempty"`
}) []struct {
	Name     string `json:"name"`
	Odds     string `json:"odds"`
	Handicap string `json:"handicap,omitempty"`
} {
	keys := make(map[string]bool)
	list := []struct {
		Name     string `json:"name"`
		Odds     string `json:"odds"`
		Handicap string `json:"handicap,omitempty"`
	}{}
	for _, entry := range selections {
		if _, value := keys[entry.Name]; !value {
			keys[entry.Name] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetDoubleChanceSelections() models.AvailableSelection {
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

func FindCorrectScoreSelection(req volleyball_models.BetEvaluationRequest) models.BetSelection {
	for _, result := range PrematchData.Results {
		for _, odd := range result.Main.Sp.CorrectSetScore.Odds {
			if odd.Header == req.Selection && odd.Name == req.ScoreLine {
				odds := fmt.Sprintf("%s", odd.Odds)
				return models.BetSelection{
					Market:    req.Market,
					Selection: req.Selection,
					Odds:      odds,
					ScoreLine: odd.Name,
				}
			}
		}
	}
	return models.BetSelection{}
}

func GetDoubleChanceOdds(selection string) string {
	// In a real implementation, you would get these from the prematch data
	// This is a simplified version with example odds
	switch selection {
	case "1X":
		return "1.10"
	case "12":
		return "1.05"
	case "X2":
		return "1.20"
	default:
		return "0"
	}
}

func GetTotalSelections() models.AvailableSelection {
	selections := []struct {
		Name     string `json:"name"`
		Odds     string `json:"odds"`
		Handicap string `json:"handicap,omitempty"`
	}{}

	for _, result := range PrematchData.Results {
		for _, odd := range result.Main.Sp.GameLines.Odds {
			if odd.Name == "Total" && odd.Header != "" {
				odds := fmt.Sprintf("%s", odd.Odds)
				selections = append(selections, struct {
					Name     string `json:"name"`
					Odds     string `json:"odds"`
					Handicap string `json:"handicap,omitempty"`
				}{
					Name:     odd.Header,
					Odds:     odds,
					Handicap: odd.Handicap,
				})
			}
		}

		for _, odd := range result.Schedule.Sp.Main {
			if odd.Name == "Total" {
				odds := fmt.Sprintf("%s", odd.Odds)
				selections = append(selections, struct {
					Name     string `json:"name"`
					Odds     string `json:"odds"`
					Handicap string `json:"handicap,omitempty"`
				}{
					Name:     "O",
					Odds:     odds,
					Handicap: odd.Handicap,
				})
				// Add under selection
				underOdds := "1.83" // Example - in real implementation you'd get this from data
				selections = append(selections, struct {
					Name     string `json:"name"`
					Odds     string `json:"odds"`
					Handicap string `json:"handicap,omitempty"`
				}{
					Name:     "U",
					Odds:     underOdds,
					Handicap: odd.Handicap,
				})
			}
		}
	}

	return models.AvailableSelection{
		Market:     "Total",
		Selections: selections,
	}
}

func GetCorrectScoreSelections() models.AvailableSelection {
	selections := []struct {
		Name     string `json:"name"`
		Odds     string `json:"odds"`
		Handicap string `json:"handicap,omitempty"`
	}{}

	for _, result := range PrematchData.Results {
		for _, odd := range result.Main.Sp.CorrectSetScore.Odds {
			odds := fmt.Sprintf("%s", odd.Odds)
			selections = append(selections, struct {
				Name     string `json:"name"`
				Odds     string `json:"odds"`
				Handicap string `json:"handicap,omitempty"`
			}{
				Name:     odd.Name,
				Odds:     odds,
				Handicap: odd.Header, // Using Header to indicate home/away
			})
		}
	}

	return models.AvailableSelection{
		Market:     "Correct Set Score",
		Selections: selections,
	}
}
