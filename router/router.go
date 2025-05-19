package router

import (
	"bet365-fiber-sim/handlers"

	_ "bet365-fiber-sim/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "/bet365 simulator service is running",
		})
	})

	app.Get("/docs/*", fiberSwagger.WrapHandler)
	api.Post("/evaluate", handlers.EvaluateCustomSelection)
	api.Get("/selections", handlers.GetAvailableSelections)
	// app.Get("/cricket/selections", handlers.GetAvailableCricketSelections)
	// app.Post("/cricket/evaluate", handlers.EvaluateCricketSelection)

}
