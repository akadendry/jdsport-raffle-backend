package jdsport_raffle_backend

import (
	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["image"]
	filename := ""

	for _, file := range files {
		filename = file.Filename

		if err := c.SaveFile(file, "./asset/products/"+filename); err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"name_file": filename,
	})
}

func UploadMobile(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["image"]
	filename := ""

	for _, file := range files {
		filename = file.Filename

		if err := c.SaveFile(file, "./asset/products/mobile/"+filename); err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"name_file": filename,
	})
}
