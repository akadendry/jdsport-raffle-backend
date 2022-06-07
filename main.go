package jdsport_raffle_backend

import (
	"log"
	"os"

	"github.com/akadendry/jdsport-raffle-backend/database"
	"github.com/akadendry/jdsport-raffle-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":" + os.Getenv("APP_PORT"))
}
