package jdsport_raffle_backend

import (
	"github.com/akadendry/jdsport-raffle-backend/v1/database"
	"github.com/akadendry/jdsport-raffle-backend/v1/models"
	"github.com/gofiber/fiber/v2"
)

func AllAccess(c *fiber.Ctx) error {
	var access []models.Access

	database.DB.Find(&access)

	return c.JSON(access)
}
