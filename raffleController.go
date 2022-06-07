package jdsport_raffle_backend

import (
	"encoding/json"
	"fmt"

	// "io/ioutil"
	// "log"
	"net/http"
	"os"
	"strings"

	"github.com/akadendry/jdsport-raffle-backend/v2/database"
	"github.com/akadendry/jdsport-raffle-backend/v2/models"
	"github.com/gosimple/slug"

	// "log"
	"time"

	// "github.com/akadendry/jdsport-raffle-backend/v1middlewares"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func AllRaffles(c *fiber.Ctx) error {
	// if err := middlewares.IsAuthorized(c, "raffles"); err != nil {
	// 	return err
	// }

	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB.Where("deleted_at IS NULL"), &models.Raffle{}, page))
}

func AddRaffle(c *fiber.Ctx) error {
	var data fiber.Map

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	ts_start_date_registration := data["start_date_registration"].(string)
	start_date_registration, _ := time.Parse("2006-01-02 15:04:05", ts_start_date_registration)

	ts_end_date_registration := data["end_date_registration"].(string)
	end_date_registration, _ := time.Parse("2006-01-02 15:04:05", ts_end_date_registration)

	ts_announcement_date := data["announcement_date"].(string)
	announcement_date, _ := time.Parse("2006-01-02 15:04:05", ts_announcement_date)

	currentTime := time.Now()
	loc, _ := time.LoadLocation("MST")
	now := time.Now().In(loc)
	fmt.Println("ZONE : ", loc, " Time : ", now) // MST
	fmt.Println("Current Time in String: ", currentTime.String())
	fmt.Println("Current Time in String: ", start_date_registration.In(loc))
	// fmt.Println(NewNullString(data["announcement_date"].(string)))
	// ts_start_date_pay := data["start_date_pay"].(string)
	// start_date_pay, _ := time.Parse("2006-01-02 15:04:05", ts_start_date_pay)

	// ts_end_date_pay := data["end_date_pay"].(string)
	// end_date_pay, _ := time.Parse("2006-01-02 15:04:05", ts_end_date_pay)

	slug := slug.Make(data["name"].(string))

	var raffleSlug models.Raffle

	var total int64
	slug_val := total
	// var slug_val int64
	// var slug_val = 0
	database.DB.Where("slug LIKE ?", "%"+slug+"%").Find(&raffleSlug).Count(&total)
	if total == 0 {
		slug_val = total
		slug = slug
	}
	if total > 0 {
		slug_val = total + 1
		slug = slug + "-" + strconv.Itoa(int(slug_val))
	}
	// var raffle models.Raffle
	raffle := models.Raffle{
		Name:                  data["name"].(string),
		StartDateRegistration: start_date_registration,
		EndDateRegistration:   end_date_registration,
		AnnouncementDate:      announcement_date,
		// StartDatePay:          time.Time{},
		// EndDatePay:            time.Time{},
		Banner:       data["banner"].(string),
		BannerMobile: data["banner_mobile"].(string),
		Copyright:    data["copyright"].(string),
		Slug:         slug,
		SlugNo:       uint(slug_val),
		CreatedAt:    time.Now(),
		CreatedBy:    data["created_by"].(string),
	}

	// database.DB.Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "id"}},                                                                                                                               // key colume
	// 	DoUpdates: clause.AssignmentColumns([]string{"name", "start_date_registration", "end_date_registration", "announcement_date", "banner", "copyright", "banner_mobile"}), // column needed to be updated
	// }).Create(&raffle)
	// database.DB.Model(&raffle).Create(map[string]interface{}{
	// 	"Name":  data["name"],
	// 	"StartDateRegistration":data["start_date_registration"],
	// 	"EndDateRegistration": data["end_date_registration"],
	// 	"AnnouncementDate": data["announcement_date"],
	// 	"Banner":data["banner"],
	// 	"BannerMobile":data["banner_mobile"],
	// 	"Copyright":data["copyright"],
	// 	"Slug": slug,
	// 	"SlugNo":slug_val,
	// 	"CreatedAt": time.Now(),
	// 	"CreatedBy": data["created_by"]})
	database.DB.Create(&raffle)

	fmt.Println(raffle.Id)

	res := Result{Code: 200, Data: raffle, Message: "Success"}
	result, err := json.Marshal(res)

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Failed",
			"data":    result,
		})
	}

	dataProdukList := data["product"].([]interface{})
	dataProduks := make([]models.GetRaffleProducts, len(dataProdukList))

	// fmt.Println(data["product"])
	// Convert map to json string
	jsonStr, err := json.Marshal(data["product"])
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	var grp []models.GetRaffleProducts
	if err := json.Unmarshal(jsonStr, &grp); err != nil {
		fmt.Println(err)
	}

	for key := range dataProduks {
		// fmt.Println("Key:", key, "=>", "Element:")
		// fmt.Println(grp[key].Name)
		// fmt.Println(grp[key].Detail)
		// dataProdukListDetail := grp[key].Detail
		// dataProduksDetail := make([]models.GetRaffleProductSizeStock, len(dataProdukListDetail))
		// fmt.Println(dataProduksDetail)
		// fmt.Println(grp[key].Detail)

		dataProduks := models.RaffleProduct{
			RaffleId:    uint(raffle.Id),
			Name:        grp[key].Name,
			Description: grp[key].Description,
			Image:       grp[key].Image,
			ImageMobile: grp[key].ImageMobile,
			CreatedAt:   time.Now(),
			CreatedBy:   data["created_by"].(string),
		}

		database.DB.Create(&dataProduks)
		for _, produksDetail := range grp[key].Detail {
			// vals := make([]string, grpd[key].RaffleProductId)
			// fmt.Println("KeyDetail :", keyd, "=>", "ElementDetail :", produksDetail)
			dataProduksSizeStock := models.RaffleProductSizeStock{
				RaffleProductId: uint(dataProduks.Id),
				Size:            produksDetail.Size,
				Stock:           produksDetail.Stock,
				UrlProduct:      produksDetail.UrlProduct,
				Sku:             produksDetail.Sku,
				CreatedAt:       time.Now(),
				CreatedBy:       data["created_by"].(string),
			}

			database.DB.Create(&dataProduksSizeStock)

		}

	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    []string{},
	})

}

func EditRaffle(c *fiber.Ctx) error {
	var data fiber.Map

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	ts_start_date_registration := data["start_date_registration"].(string)
	start_date_registration, _ := time.Parse("2006-01-02 15:04:05", ts_start_date_registration)

	ts_end_date_registration := data["end_date_registration"].(string)
	end_date_registration, _ := time.Parse("2006-01-02 15:04:05", ts_end_date_registration)

	ts_announcement_date := data["announcement_date"].(string)
	announcement_date, _ := time.Parse("2006-01-02 15:04:05", ts_announcement_date)

	// ts_start_date_pay := data["start_date_pay"].(string)
	// start_date_pay, _ := time.Parse("2006-01-02 15:04:05", ts_start_date_pay)

	// ts_end_date_pay := data["end_date_pay"].(string)
	// end_date_pay, _ := time.Parse("2006-01-02 15:04:05", ts_end_date_pay)

	raffle_id, _ := strconv.Atoi(data["raffle_id"].(string))
	slug := slug.Make(data["name"].(string))

	var raffleSlug models.Raffle

	var total int64
	slug_val := total
	// var slug_val int64
	// var slug_val = 0
	database.DB.Where("id = ?", raffle_id).Where("slug LIKE ?", "%"+slug+"%").Find(&raffleSlug).Count(&total)
	// fmt.Println(total)
	if total == 0 {
		slug_val = total
		slug = slug
	}
	if total > 0 {
		slug_val = total + 1
		slug = slug + "-" + strconv.Itoa(int(slug_val))
	}

	raffle := models.Raffle{
		Id:                    uint(raffle_id),
		Name:                  data["name"].(string),
		StartDateRegistration: start_date_registration,
		EndDateRegistration:   end_date_registration,
		AnnouncementDate:      announcement_date,
		// StartDatePay:          start_date_pay,
		// EndDatePay:            end_date_pay,
		Banner:       data["banner"].(string),
		BannerMobile: data["banner_mobile"].(string),
		Copyright:    data["copyright"].(string),
		Slug:         slug,
		SlugNo:       uint(slug_val),
		UpdatedAt:    time.Now(),
		UpdatedBy:    data["updated_by"].(string),
	}

	database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},                                                                                                                                                           // key colume
		DoUpdates: clause.AssignmentColumns([]string{"name", "start_date_registration", "end_date_registration", "announcement_date", "banner", "copyright", "banner_mobile", "updated_at", "updated_by"}), // column needed to be updated
	}).Create(&raffle)

	// database.DB.Create(&raffle)

	// fmt.Println(raffle.Id)

	res := Result{Code: 200, Data: raffle, Message: "Success"}
	result, err := json.Marshal(res)

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Failed",
			"data":    result,
		})
	}

	dataProdukList := data["product"].([]interface{})
	dataProduks := make([]models.GetRaffleProducts, len(dataProdukList))

	// fmt.Println(data["product"])
	// Convert map to json string
	jsonStr, err := json.Marshal(data["product"])
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	var grp []models.GetRaffleProducts
	if err := json.Unmarshal(jsonStr, &grp); err != nil {
		fmt.Println(err)
	}

	for key := range dataProduks {
		produks_id, _ := strconv.Atoi(grp[key].Id)

		inDataProduks := models.RaffleProduct{
			Id:          uint(produks_id),
			RaffleId:    uint(raffle.Id),
			Name:        grp[key].Name,
			Description: grp[key].Description,
			Image:       grp[key].Image,
			ImageMobile: grp[key].ImageMobile,
			UpdatedAt:   time.Now(),
			UpdatedBy:   data["updated_by"].(string),
		}

		// database.DB.Create(&inDataProduks)
		database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},                                                                                               // key colume
			DoUpdates: clause.AssignmentColumns([]string{"name", "raffle_id", "description", "image", "image_mobile", "updated_at", "updated_by"}), // column needed to be updated
		}).Create(&inDataProduks)

		for _, produksDetail := range grp[key].Detail {
			// vals := make([]string, grpd[key].RaffleProductId)
			// fmt.Println("KeyDetail :", keyd, "=>", "ElementDetail :", produksDetail)
			produks_detail_id, _ := strconv.Atoi(produksDetail.Id)

			dataProduksSizeStock := models.RaffleProductSizeStock{
				Id:              uint(produks_detail_id),
				RaffleProductId: uint(inDataProduks.Id),
				Size:            produksDetail.Size,
				Stock:           produksDetail.Stock,
				UrlProduct:      produksDetail.UrlProduct,
				Sku:             produksDetail.Sku,
				UpdatedAt:       time.Now(),
				UpdatedBy:       data["updated_by"].(string),
			}

			// database.DB.Create(&dataProduksSizeStock)
			database.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},                                                                                              // key colume
				DoUpdates: clause.AssignmentColumns([]string{"raffle_product_id", "size", "stock", "url_product", "sku", "updated_at", "updated_by"}), // column needed to be updated
			}).Create(&dataProduksSizeStock)

		}

	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    []string{},
	})

}

// func AllRaffles(c *fiber.Ctx) error {
// 	var raffles []models.Raffle

// 	database.DB.Find(&raffles)

// 	return c.JSON(raffles)
// }

// func CreateRaffle(c *fiber.Ctx) error {
// 	if err := middlewares.IsAuthorized(c, "raffle"); err != nil {
// 		return err
// 	}

// 	var data map[string]string

// 	if err := c.BodyParser(&data); err != nil {
// 		return err
// 	}

// 	if data["password"] != data["password_confirm"] {
// 		c.Status(400)
// 		return c.JSON(fiber.Map{
// 			"message": "passwords do not match",
// 		})
// 	}
// 	role_id, _ := strconv.Atoi(data["role_id"])

// 	raffle := models.Raffle{
// 		FirstName: data["first_name"],
// 		LastName:  data["last_name"],
// 		Email:     data["email"],
// 		RoleId:    uint(role_id),
// 	}

// 	raffle.SetPassword(data["password"])

// 	database.DB.Create(&raffle)

// 	return c.JSON(raffle)
// }

func GetRaffleById(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id := data["id"]

	var result models.Raffle_New
	database.DB.Raw("SELECT * FROM raffles WHERE id = ?", id).Where("deleted_at IS NULL").Scan(&result)

	if result.Id == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Data Not Found",
			"data":    []string{},
		})
	}

	current_time := time.Now()

	if current_time.After(result.EndDateRegistration) {
		return c.JSON(fiber.Map{
			"status":  "202",
			"message": "Raffle sudah berakhir!",
			"data":    result,
		})
	} else {
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Success",
			"data":    result,
		})
	}
}

func GetRaffleBySlug(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	slug := data["slug"]

	var result models.Raffle_New
	database.DB.Raw("SELECT * FROM raffles WHERE slug = ?", slug).Where("deleted_at IS NULL").Scan(&result)

	if result.Id == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Data Not Found",
			"data":    []string{},
		})
	}

	current_time := time.Now()

	if current_time.After(result.EndDateRegistration) {
		return c.JSON(fiber.Map{
			"status":  "202",
			"message": "Raffle sudah berakhir!",
			"data":    result,
		})
	} else {
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Success",
			"data":    result,
		})
	}
}

func UpdateRaffle(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	raffle := models.Raffle{
		Id: uint(id),
	}

	if err := c.BodyParser(&raffle); err != nil {
		return err
	}

	database.DB.Model(&raffle).Updates(raffle)

	return c.JSON(raffle)
}

func GetRaffleByFilter(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var result []models.Raffle
	var data_title string
	var data_status string
	var data_start_date_registration string
	var data_end_date_registration string

	// tx := database.DB

	if data["title"] != "" {
		data_title = "raffles.name LIKE  '%" + data["title"] + "%'"
	} else {
		data_title = ""
	}

	if data["status"] != "" {
		data_status = "raffles.status =  " + data["status"]
	} else {
		data_status = ""
	}

	if data["start_date_registration"] != "" {
		data_start_date_registration = "raffles.start_date_registration =  '" + data["start_date_registration"] + "'"
	} else {
		data_start_date_registration = ""
	}

	if data["end_date_registration"] != "" {
		data_end_date_registration = "raffles.end_date_registration = '" + data["end_date_registration"] + "'"
	} else {
		data_end_date_registration = ""
	}

	tx := database.DB.Where("deleted_at IS NULL").
		Where(data_title).
		Where(data_status).
		Where(data_start_date_registration).
		Where(data_end_date_registration).
		Find(&result)
	var total int64
	tx.Find(&result).Count(&total)
	fmt.Println("Total  :", total)
	if total == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Data Not Found",
			"data":    []string{},
		})
	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    result,
	})
}

func DeleteRaffle(c *fiber.Ctx) error {
	raffle_id, _ := strconv.Atoi(c.Params("id"))
	fmt.Println(raffle_id)
	raffle := models.Raffle{
		Id:        uint(raffle_id),
		DeletedAt: time.Now(),
		DeletedBy: c.Params("deleted_by"),
	}
	database.DB.Updates(&raffle)

	var raffleProduct []models.RaffleProduct
	database.DB.Where("raffle_id = ?", raffle_id).Find(&raffleProduct)

	// fmt.Println(raffleProduct[0])

	for _, value := range raffleProduct {
		// fmt.Println("index : ", key, " value : ", value.Id)
		inDataProduks := models.RaffleProduct{
			Id:        uint(value.Id),
			RaffleId:  uint(raffle_id),
			DeletedAt: time.Now(),
			DeletedBy: c.Params("deleted_by"),
		}

		database.DB.Updates(&inDataProduks)

		var raffleProductDetail []models.RaffleProductSizeStock
		database.DB.Where("raffle_product_id = ?", value.Id).Find(&raffleProductDetail)
		for _, valueDetail := range raffleProductDetail {
			// fmt.Println("index : ", key, " value : ", valueDetail.Id)
			inDataProduksDetail := models.RaffleProductSizeStock{
				Id:        uint(valueDetail.Id),
				DeletedAt: time.Now(),
				DeletedBy: c.Params("deleted_by"),
			}

			database.DB.Updates(&inDataProduksDetail)
		}

	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Delete Success",
		"data":    []string{},
	})

}

//Ini untuk CMS
func GetDetailRaffleByIdNew(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id := data["id"]

	var result models.Raffle
	database.DB.Raw("SELECT * FROM raffles WHERE id = ?", id).Where("deleted_at IS NULL").
		Scan(&result)

	if result.Id == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Data Not Found",
			"data":    []string{},
		})
	}
	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    result,
	})
}

func GetDetailRaffleById(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	raffle_id, _ := strconv.Atoi(data["id"])

	var result []models.Raffle

	database.DB.Table("raffles").Where("id", raffle_id).Preload("RaffleProduct").Preload("RaffleProduct.RaffleProductSizeStock").Find(&result)

	// database.DB.Raw("SELECT * FROM raffles WHERE id = ?", raffle_id).
	// 	Where("raffles.deleted_at IS NULL").
	// 	Joins("JOIN raffle_products on raffles.id = raffle_products.raffle_id").
	// 	Joins("JOIN raffle_product_size_stocks on raffle_products.id = raffle_product_size_stocks.raffle_product_id").
	// 	Preload("RaffleProduct").
	// 	Preload("RaffleProductSizeStock").
	// 	Scan(&result)

	// tx := database.DB

	// tx = database.DB.Where("raffle_products.raffle_id = ?", raffle_id).
	// 	Where("raffle_products.deleted_at IS NULL").
	// 	Joins("JOIN raffles on raffles.id = raffle_products.raffle_id").
	// 	Joins("JOIN raffle_product_size_stocks on raffle_products.id = raffle_product_size_stocks.raffle_product_id").
	// 	Preload("Raffle").
	// 	Preload("RaffleProductSizeStock").
	// 	Find(&raffle)
	// fmt.Println(tx)
	// for _, ar := range participant {
	// 	fmt.Println(ar.User.FirstName)
	// }

	// fmt.Println(asd)
	// 	tx.Find(&participant)
	// var total int64
	// tx.Find(&raffle).Count(&total)

	// if total == 0 {
	// 	// c.Status(404)
	// 	return c.JSON(fiber.Map{
	// 		"status":  "404",
	// 		"message": "Data Not Found",
	// 		"data":    []string{},
	// 	})
	// }

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    result,
	})
}

func EditRaffleStatus(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.Atoi(data["raffle_id"])

	if data["status"] == "3" {

		fmt.Println("raffle_id", id)
		status := os.Getenv("STATUS")
		var link_service_email string
		if status == "dev" {
			link_service_email = os.Getenv("EMAIL_SERVICE_DEV")
		} else {
			link_service_email = os.Getenv("EMAIL_SERVICE_PROD")
		}
		// fmt.Println("link_service_email", link_service_email)

		var raffle models.Raffle
		database.DB.Where("id = ?", id).Find(&raffle)
		const (
			DDMMYYYYhhmmss = "02-01-2006 15:04:05"
		)

		// var raffleCek models.Raffle
		// var total int64

		// database.DB.Where("id = ?", id).Where("is_email = ?", 1).Find(&raffleCek).Count(&total)
		// fmt.Println("total rafle : ", total)

		// if total == 0 {

		var participant []models.Participant
		database.DB.Where("participants.raffle_id = ?", id).
			Where("participants.deleted_at IS NULL").
			Where("participants.status = ?", 2).
			Where("participants.is_email = ?", 0).
			Joins("JOIN users on participants.user_id = users.id").
			Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
			Joins("JOIN raffle_products on participants.raffle_product_id=raffle_products.id").
			Preload("User").
			Preload("RaffleProduct").
			Preload("RaffleProductSizeStock").
			Find(&participant)

		// fmt.Println("data Participant :", participant)
		for _, dataParticipant := range participant {
			fmt.Println("data Participant :", dataParticipant.User.Email)
			// fmt.Println("data Participant :", dataParticipant.RaffleProductSizeStock.UrlProduct)
			// fmt.Println("dataParticipant.Id :", dataParticipant.Id)
			ts_start_date_pay := dataParticipant.StartDatePay
			// fmt.Println("data start_date_pay :", ts_start_date_pay.Format(DDMMYYYYhhmmss))

			start_date_pay := ts_start_date_pay.Format(DDMMYYYYhhmmss)
			fmt.Println("data end_date_pay :", start_date_pay)

			ts_end_date_pay := dataParticipant.EndDatePay
			// fmt.Println("data end_date_pay :", ts_end_date_pay.Format(DDMMYYYYhhmmss ))
			end_date_pay := ts_end_date_pay.Format(DDMMYYYYhhmmss)
			fmt.Println("data end_date_pay :", end_date_pay)
			var participantUpdate models.Participant
			database.DB.Where("id = ?", dataParticipant.Id).Model(&participantUpdate).Updates(map[string]interface{}{"IsEmail": 1, "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})
			payload := strings.NewReader(`{` + "" + `
				"email": "` + dataParticipant.User.Email + `", 
				"host": "jdsport",` + "" + `
				"title": "Pengumuman Raffle",` + "" + `
				"body":  "<h2>Mohon Maaf ` + dataParticipant.User.FirstName + ` - ` + dataParticipant.User.LastName + ` Kamu belum beruntung untuk Raffle ` + raffle.Name + `!</h2>` + "" +
				`<h3> Terimakasih Atas Partisipasinya!</h3>"` + "" + `
			}`)
			// fmt.Println("LINK:", payload)
			service_mail, err := http.NewRequest("POST", link_service_email, payload)
			service_mail.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			responsed, error := client.Do(service_mail)
			if error != nil {
				panic(error)
			}
			defer responsed.Body.Close()
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			// fmt.Println("responsed Body:", responsed.Body)
			// responseData, err := ioutil.ReadAll(responsed.Body)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Println(string(responseData))
			// fmt.Println("serviceEmailResponseData:", string(responseData))
			// var responseObject models.ResponseServiceEmail
			// json.Unmarshal(responseData, &responseObject)

		}

		var raffleNow models.Raffle
		database.DB.Where("id = ?", id).Model(&raffleNow).Updates(map[string]interface{}{"Status": data["status"], "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Susses",
			"data":    "[]",
		})

	}
	var raffleNow models.Raffle
	database.DB.Where("id = ?", id).Model(&raffleNow).Updates(map[string]interface{}{"Status": data["status"], "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    []string{},
	})
}

func GetTotalWinerParticipanAndTotalStock(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := data["raffle_id"]
	var raffleProductResume []models.RaffleProductResume

	sql := "SELECT " +
		"raffle_products.`name` AS name_product, " +
		"raffle_product_size_stocks.size, " +
		"raffle_product_size_stocks.stock, " +
		"( " +
		"SELECT COUNT(*) " +
		"FROM participants " +
		"WHERE participants.raffle_product_id = raffle_products.id " +
		"AND participants.raffle_id = " + id + " " +
		"AND participants.raffle_product_size_stock_id = raffle_product_size_stocks.id " +
		"AND participants.`status` = 1 " +
		") AS jumlah_pemenang " +
		"FROM " +
		"raffle_products " +
		"LEFT JOIN raffle_product_size_stocks ON raffle_product_size_stocks.raffle_product_id = raffle_products.id " +
		"WHERE raffle_products.raffle_id = " + id
	tx := database.DB.Debug().Raw(sql).Scan(&raffleProductResume)

	var total int64
	tx.Find(&raffleProductResume).Count(&total)
	if total == 0 || id == "" {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Data Not Found",
			"data":    []string{},
		})
	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    raffleProductResume,
	})

}

func GetSelectOptionRaffle(c *fiber.Ctx) error {

	var raffle []models.RaffleProductSize

	sql := "SELECT " +
		"raffle_product_size_stocks.size " +
		"FROM " +
		"raffles " +
		"LEFT JOIN raffle_products ON raffle_products.raffle_id = raffles.id " +
		"LEFT JOIN raffle_product_size_stocks ON raffle_product_size_stocks.raffle_product_id = raffle_products.id " +
		"GROUP BY raffle_product_size_stocks.size " +
		"ORDER BY raffle_product_size_stocks.size ASC "
	database.DB.Debug().Raw(sql).Scan(&raffle)

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    raffle,
	})

}
