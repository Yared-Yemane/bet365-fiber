package main

import (
	"bet365-fiber-sim/router"
	"bet365-fiber-sim/utils"
	cricket_utils "bet365-fiber-sim/utils/cricket"
	volleyball_utils "bet365-fiber-sim/utils/volleyball"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// @title Betting Evaluation API
// @version 1.0
// @description API for evaluating betting selections against match results
// @contact.name API Support
// @contact.email yaredyemane1@gmail.com
// @license.name Yared Yemane
// @host localhost:8080
// @BasePath /api/v1/
func main() {
	app := fiber.New()

	volleyball_utils.InitHandlers(app)
	cricket_utils.InitCricketHandlers()
	utils.ConfigCORS(app)

	router.SetupRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("INTERNAL_PORT"))))
}
