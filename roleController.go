package jdsport_raffle_backend

import (
	"strconv"

	"github.com/akadendry/jdsport-raffle-backend/v2/database"
	"github.com/akadendry/jdsport-raffle-backend/v2/models"
	"github.com/gofiber/fiber/v2"
)

func AllRoles(c *fiber.Ctx) error {
	var roles []models.Role

	database.DB.Find(&roles)

	return c.JSON(roles)
}

func CreateRole(c *fiber.Ctx) error {
	var roleDto fiber.Map

	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	list := roleDto["access"].([]interface{})

	accesses := make([]models.Access, len(list))

	for i, accessId := range list {
		id, _ := strconv.Atoi(accessId.(string))

		accesses[i] = models.Access{
			Id: uint(id),
		}
	}

	role := models.Role{
		Name:   roleDto["name"].(string),
		Access: accesses,
	}

	database.DB.Create(&role)

	return c.JSON(role)
}

func GetRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	role := models.Role{
		Id: uint(id),
	}

	database.DB.Preload("Access").Find(&role)

	return c.JSON(role)
}

func UpdateRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var roleDto fiber.Map

	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	list := roleDto["access"].([]interface{})

	accesses := make([]models.Access, len(list))

	for i, permissionId := range list {
		id, _ := permissionId.(float64)

		accesses[i] = models.Access{
			Id: uint(id),
		}
	}

	var result interface{}

	database.DB.Table("role_accesses").Where("role_id", id).Delete(&result)

	role := models.Role{
		Id:     uint(id),
		Name:   roleDto["name"].(string),
		Access: accesses,
	}

	database.DB.Model(&role).Updates(role)

	return c.JSON(role)
}

func DeleteRole(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	role := models.Role{
		Id: uint(id),
	}

	database.DB.Delete(&role)

	return nil
}
