package jdsport_raffle_backend

import (
	"github.com/akadendry/jdsport-raffle-backend/v1/database"
	"github.com/akadendry/jdsport-raffle-backend/v1/models"
	"github.com/gofiber/fiber/v2"
)

func GetRaffleProductByRaffleId(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var raffleProduct []models.RaffleProduct
	raffleId := data["raffle_id"]
	database.DB.Where("raffle_id = ?", raffleId).Where("deleted_at IS NULL").Find(&raffleProduct)

	var raffleProductCount models.RaffleProduct
	raffleIdcount := data["raffle_id"]
	var total int64
	database.DB.Where("raffle_id = ?", raffleIdcount).Where("deleted_at IS NULL").Find(&raffleProductCount).Count(&total)

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    raffleProduct,
		"total":   total,
	})
}

func GetRaffleProductById(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var raffleProduct models.RaffleProduct
	idProduct := data["product_id"]
	database.DB.Where("id = ?", idProduct).Where("deleted_at IS NULL").Find(&raffleProduct)

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    raffleProduct,
	})
}
