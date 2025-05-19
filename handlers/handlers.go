package handlers

import (
	models "bet365-fiber-sim/models"
	cricket_models "bet365-fiber-sim/models/cricket"
	volleyball_models "bet365-fiber-sim/models/volleyball"
	cricket_utils "bet365-fiber-sim/utils/cricket"
	volleyball_utils "bet365-fiber-sim/utils/volleyball"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get available betting selections
// @Description Retrieves all available betting markets and selections from prematch data
// @Tags Selections
// @Accept  json
// @Produce  json
// @Success 200 {array} models.AvailableSelection "List of available selections grouped by market"
// @Failure 404 {object} object "No prematch data available"
// @Router /selections [get]
func GetAvailableSelections(c *fiber.Ctx) error {

	var available []models.AvailableSelection
	sport_type := c.Query("sport_type")

	if sport_type == "volleyball" {
		if len(volleyball_utils.PrematchData.Results) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No prematch data available",
			})
		}

		available = []models.AvailableSelection{
			volleyball_utils.Get1X2Selections(),
			volleyball_utils.GetTotalSelections(),
			volleyball_utils.GetCorrectScoreSelections(),
			volleyball_utils.GetDoubleChanceSelections(),
		}

		// return c.JSON(available)
	} else if sport_type == "cricket" {
		if len(cricket_utils.PrematchData.Results) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No prematch data available",
			})
		}

		availableTemp := []models.AvailableSelection{
			cricket_utils.GetCricket1X2Selections(),
			cricket_utils.GetCricketTotalRunsSelections(),
			cricket_utils.GetCricketDoubleChanceSelections(),
		}

		available = availableTemp

	}
	return c.JSON(available)
}

// @Summary Evaluate a betting selection
// @Description Evaluates a specific betting selection against the match results
// @Tags Evaluation
// @Accept json
// @Produce json
// @Param request body models.BetEvaluationRequest true "Bet selection to evaluate"
// @Success 200 {object} models.EvaluationResult "Evaluation result with outcome"
// @Failure 400 {object} object "Invalid request body or parameters"
// @Failure 404 {object} object "No result data available"
// @Router /evaluate [post]
func EvaluateCustomSelection(c *fiber.Ctx) error {
	var req volleyball_models.BetEvaluationRequest
	var result models.EvaluationResult
	sport_type := c.Query("sport_type")

	if sport_type == "volleyball" {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if len(volleyball_utils.ResultData.Results) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No result data available",
			})
		}

		selection := volleyball_utils.CreateSelectionFromRequest(req)
		if selection.Market == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid selection parameters",
			})
		}

		result = volleyball_utils.EvaluateSelection(selection, volleyball_utils.ResultData)
	} else if sport_type == "cricket" {
		var req cricket_models.BetEvaluationRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if len(cricket_utils.ResultData.Results) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No result data available",
			})
		}

		selection := cricket_utils.CreateCricketSelectionFromPrematch(cricket_utils.PrematchData, req.Market, req.Selection, req.Handicap)
		if selection.Market == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid selection parameters",
			})
		}

		result = cricket_utils.EvaluateCricketSelection(selection, cricket_utils.ResultData)
	}

	return c.JSON(result)
}

// @Summary Get available cricket betting selections
// @Description Retrieves all available cricket betting markets and selections from prematch data
// @Tags Cricket Selections
// @Accept  json
// @Produce  json
// @Success 200 {array} models.AvailableSelection "List of available selections grouped by market"
// @Failure 404 {object} object "No prematch data available"
// @Router /cricket/selections [get]
// func GetAvailableCricketSelections(c *fiber.Ctx) error {
// 	if len(cricket_models.PrematchResponse.Results) == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "No prematch data available",
// 		})
// 	}

// 	available := []models.AvailableSelection{
// 		cricket_utils.GetCricket1X2Selections(),
// 		cricket_utils.GetCricketTotalRunsSelections(),
// 		cricket_utils.GetCricketDoubleChanceSelections(),
// 	}

// 	return c.JSON(available)
// }

// @Summary Evaluate a cricket betting selection
// @Description Evaluates a specific cricket betting selection against the match results
// @Tags Cricket Evaluation
// @Accept json
// @Produce json
// @Param request body models.BetEvaluationRequest true "Bet selection to evaluate"
// @Success 200 {object} models.EvaluationResult "Evaluation result with outcome"
// @Failure 400 {object} object "Invalid request body or parameters"
// @Failure 404 {object} object "No result data available"
// @Router /cricket/evaluate [post]
// func EvaluateCricketSelection(c *fiber.Ctx) error {

// }

// func EvaluateBet(c *fiber.Ctx) error {
// 	// Read prematch data from file
// 	volleyball_utils.PrematchData, err := volleyball_utils.ReadPrematchData("data/volleyball_prematch.json")
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": fmt.Sprintf("Failed to read prematch data: %v", err),
// 		})
// 	}

// 	// Read result data from file
// 	volleyball_utils.ResultData, err := volleyball_utils.Readvolleyball_utils.ResultData("data/volleyball_result.json")
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": fmt.Sprintf("Failed to read result data: %v", err),
// 		})
// 	}

// 	// Sample selections from prematch data
// 	selections := []models.BetSelection{
// 		// 1X2 Market (Win/Draw/Win)
// 		volleyball_utils.CreateSelectionFromPrematch(volleyball_utils.PrematchData, "Winner", "1"), // Home win
// 		volleyball_utils.CreateSelectionFromPrematch(volleyball_utils.PrematchData, "Winner", "X"), // Draw (if available)
// 		volleyball_utils.CreateSelectionFromPrematch(volleyball_utils.PrematchData, "Winner", "2"), // Away win

// 		// Over/Under
// 		volleyball_utils.CreateSelectionFromPrematch(volleyball_utils.PrematchData, "Total", "1", "O 177.5"), // Over
// 		volleyball_utils.CreateSelectionFromPrematch(volleyball_utils.PrematchData, "Total", "2", "U 177.5"), // Under

// 		// Correct Score
// 		volleyball_utils.CreateCorrectScoreSelection(volleyball_utils.PrematchData, "1", "3-0"), // Home 3-0
// 		volleyball_utils.CreateCorrectScoreSelection(volleyball_utils.PrematchData, "2", "0-3"), // Away 0-3

// 		// Double Chance
// 		volleyball_utils.CreateDoubleChanceSelection(volleyball_utils.PrematchData, "1X"), // Home or Draw
// 		volleyball_utils.CreateDoubleChanceSelection(volleyball_utils.PrematchData, "12"), // Home or Away
// 		volleyball_utils.CreateDoubleChanceSelection(volleyball_utils.PrematchData, "X2"), // Draw or Away
// 	}

// 	// Evaluate each selection
// 	results := make([]models.EvaluationResult, 0)
// 	for _, selection := range selections {
// 		result := volleyball_utils.EvaluateSelection(selection, volleyball_utils.ResultData)
// 		results = append(results, result)
// 	}

// 	return c.JSON(results)
// }
