package jdsport_raffle_backend

import (
	"os"

	// "log"
	"time"

	// "github.com/akadendry/jdsport-raffle-backend/v1middlewares"

	"strconv"

	"github.com/akadendry/jdsport-raffle-backend/v1/database"
	"github.com/akadendry/jdsport-raffle-backend/v1/models"
	"github.com/gofiber/fiber/v2"
)

//untuk api team backend jdsport
//untuk get raffle yang hari ini
func GetRaffleByDateNow(c *fiber.Ctx) error {
	ua := c.Request().Header.Peek("User-Agent")
	uas := string(ua)

	pass := c.Request().Header.Peek("Password")
	passs := string(pass)

	user_agent := os.Getenv("USER_AGENT")
	password := os.Getenv("PASSWORD")

	if uas == user_agent && passs == password {
		current_time := time.Now()

		var result []models.Raffle
		database.DB.Table("raffles").Where("end_date_registration > ?", current_time).Where("deleted_at is NULL").Scan(&result)

		var total int64
		database.DB.Table("raffles").Select("COUNT(id)").Where("end_date_registration > ?", current_time).Where("deleted_at is NULL").Scan(&total)

		if total < 1 {
			// c.Status(404)
			return c.JSON(fiber.Map{
				"status":  "202",
				"message": "Data Not Found",
				"data":    []string{},
			})
		}
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Success",
			"data":    result,
		})
	} else {
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Anda tidak diperkenankan mengakses endpoint ini!",
			"data":    []string{},
		})
	}

}

//untuk cek ketika checkout dan pembayaran
func CheckUserWinner(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	ua := c.Request().Header.Peek("User-Agent")
	uas := string(ua)

	pass := c.Request().Header.Peek("Password")
	passs := string(pass)

	user_agent := os.Getenv("USER_AGENT")
	password := os.Getenv("PASSWORD")

	if uas == user_agent && passs == password {
		erajaya_club_user_id, _ := data["user_id"]
		sku, _ := data["sku"]
		current_time := time.Now()
		quantity, _ := strconv.Atoi(data["quantity"])

		if quantity > 1 {
			return c.JSON(fiber.Map{
				"status":  "202",
				"message": "Quantity produk melebihi jumlah yang ditetapkan",
				"data":    []string{},
			})
		}

		var total int
		database.DB.Table("participants").
			Select("COUNT(participants.id)").
			Where("participants.erajaya_club_user_id = ?", erajaya_club_user_id).
			Where("participants.deleted_at IS NULL").
			Where("participants.status = ?", 1).
			Where("participants.status_transaction = ?", 0).
			Where("raffle_product_size_stocks.sku = ?", sku).
			Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id = raffle_product_size_stocks.id").
			Scan(&total)

		if total == 1 {
			var total2 int
			database.DB.Table("participants").
				Select("COUNT(participants.id)").
				Where("participants.erajaya_club_user_id = ?", erajaya_club_user_id).
				Where("participants.deleted_at IS NULL").
				Where("participants.status = ?", 1).
				Where("participants.status_transaction = ?", 0).
				Where("participants.end_date_pay > ?", current_time).
				Where("raffle_product_size_stocks.sku = ?", sku).
				Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id = raffle_product_size_stocks.id").
				Scan(&total2)
			if total2 == 1 {
				return c.JSON(fiber.Map{
					"status":  "200",
					"message": "Success",
					"data":    []string{},
				})
			} else {
				return c.JSON(fiber.Map{
					"status":  "202",
					"message": "Batas waktu pembayaran sudah lewat",
					"data":    []string{},
				})
			}

		} else {
			return c.JSON(fiber.Map{
				"status":  "202",
				"message": "User pemenang tidak ditemukan",
				"data":    []string{},
			})
		}
	} else {
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Anda tidak diperkenankan mengakses endpoint ini!",
			"data":    []string{},
		})
	}
}

func UpdateTransactionStatus(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	ua := c.Request().Header.Peek("User-Agent")
	uas := string(ua)

	pass := c.Request().Header.Peek("Password")
	passs := string(pass)

	user_agent := os.Getenv("USER_AGENT")
	password := os.Getenv("PASSWORD")

	if uas == user_agent && passs == password {
		erajaya_club_user_id, _ := data["user_id"]
		sku, _ := data["sku"]
		current_time := time.Now()
		status, _ := strconv.Atoi(data["status"])

		//get data raffle product id
		var raffle_product_id uint
		database.DB.Table("raffle_product_size_stocks").
			Select("raffle_product_id").
			Where("deleted_at IS NULL").
			Where("sku= ?", sku).
			Scan(&raffle_product_id)

		// update status
		database.DB.Table("participants").
			Where("erajaya_club_user_id = ?", erajaya_club_user_id).
			Where("raffle_product_id = ?", raffle_product_id).
			Updates(map[string]interface{}{"status_transaction": status, "updated_at": current_time, "updated_by": "Hit By BackEnd JD"})

		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Success",
			"data":    []string{},
		})
	} else {
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Anda tidak diperkenankan mengakses endpoint ini!",
			"data":    []string{},
		})
	}
}
