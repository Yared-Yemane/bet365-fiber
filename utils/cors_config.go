package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ConfigCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Add your frontend URLs
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "*",
	}))
}
