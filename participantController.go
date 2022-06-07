package jdsport_raffle_backend

import (
	"encoding/json"
	"fmt"

	// "io/ioutil"
	// "log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akadendry/jdsport-raffle-backend/v2/database"
	"github.com/akadendry/jdsport-raffle-backend/v2/models"
	"github.com/gofiber/fiber/v2"
)

func CreateParticipant(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	raffleId, _ := strconv.Atoi(data["raffle_id"])
	userId, _ := strconv.Atoi(data["user_id"])
	idErajayaClubUser, _ := strconv.Atoi(data["erajaya_club_user_id"])
	raffleProductId, _ := strconv.Atoi(data["raffle_product_id"])
	raffleProductSizeStockId, _ := strconv.Atoi(data["raffle_product_size_stock_id"])

	var participantTotal models.Participant
	raffle_id, _ := strconv.Atoi(data["raffle_id"])
	var totalparticipant int64
	database.DB.Find(&participantTotal, "raffle_id = ? ", raffle_id).Count(&totalparticipant)

	var total int64

	participant := models.Participant{
		RaffleId:                 uint(raffleId),
		UserId:                   uint(userId),
		ErajayaClubUserId:        uint(idErajayaClubUser),
		RaffleProductId:          uint(raffleProductId),
		RaffleProductSizeStockId: uint(raffleProductSizeStockId),
		QueueNo:                  uint(totalparticipant + 1),
		CreatedAt:                time.Now(),
		CreatedBy:                data["erajaya_club_user_id"],
	}

	database.DB.Find(&participant, "raffle_id = ? AND user_id  = ? AND erajaya_club_user_id = ? ", raffleId, userId, idErajayaClubUser).Count(&total)

	if total >= 1 {
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Success",
			"data":    []string{},
		})
	}

	database.DB.Create(&participant)

	res := Result{Code: 200, Data: participant, Message: "Success"}
	result, err := json.Marshal(res)

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Failed",
			"data":    result,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    []string{},
	})

}

func CekUserParticipant(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var participant []models.Participant
	raffle_id, _ := strconv.Atoi(data["raffle_id"])
	erajaya_club_user_id, _ := strconv.Atoi(data["erajaya_club_user_id"])
	var total int64
	database.DB.Find(&participant, "erajaya_club_user_id = ? AND raffle_id = ? ", erajaya_club_user_id, raffle_id).Count(&total)

	if total >= 1 {
		return c.JSON(fiber.Map{
			"data":    "[]",
			"message": "Anda sudah mengikuti raffle ini, Silahkan tunggu hasil dari raffle ini!",
			"status":  "400",
		})
	}

	return c.JSON(fiber.Map{
		"data":    "[]",
		"message": "Success",
		"status":  "200",
	})

}

func RaffleParticipant(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id, _ := strconv.Atoi(data["raffle_id"])

	var participant []models.Participant

	// tx := database.DB

	tx := database.DB.Where("participants.raffle_id = ?", id).
		Where("participants.deleted_at IS NULL").
		Joins("JOIN users on participants.user_id = users.id").
		Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
		Joins("JOIN raffle_products on participants.raffle_product_id=raffle_products.id").
		Preload("User").
		Preload("RaffleProduct").
		Preload("RaffleProductSizeStock").
		Find(&participant)
	// fmt.Println(tx)
	// for _, ar := range participant {
	// 	fmt.Println(ar.User.FirstName)
	// }

	// fmt.Println(participant)
	// 	tx.Find(&participant)
	var total int64
	tx.Find(&participant).Count(&total)

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
		"data":    participant,
	})
}

func GenerateWinerParticipant(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id, _ := strconv.Atoi(data["raffle_id"])

	var participantStock models.Participant
	database.DB.Where("participants.raffle_id = ?", id).
		Where("users.status = ?", 1).
		Where("participants.deleted_at IS NULL").
		Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
		Joins("JOIN users on participants.user_id = users.id").
		Preload("RaffleProductSizeStock").
		Find(&participantStock)
	// fmt.Println(participantStock.RaffleProductSizeStock.Stock)
	totalWiner, _ := strconv.Atoi(participantStock.RaffleProductSizeStock.Stock)

	var raffle models.Raffle
	database.DB.Where("id = ?", id).Model(&raffle).Updates(map[string]interface{}{"Status": 0, "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})

	var participantTotal models.Participant
	var totalParticipant int64
	database.DB.Where("participants.raffle_id = ?", id).Find(&participantTotal).Count(&totalParticipant)

	var participant models.Participant
	// fmt.Println("168 totalParticipant : ",totalParticipant)
	// fmt.Println("173 fortotalWinerMax : ",int(totalWiner))
	// database.DB.Where("raffle_id = ?", id).Model(&participant).Updates(map[string]interface{}{"Status": 0, "UpdatedAt": time.Now(), "UpdatedBy":"Generate By System"})
	database.DB.Where("raffle_id = ?", id).Model(&participant).Updates(map[string]interface{}{"Status": 2, "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})
	// min := 1
	// max := int(totalParticipant+1)
	// noWiner := rand.Intn(max - min) + min
	rand.Seed(time.Now().Unix())
	min := 1
	max := int(totalParticipant + 1)
	permutation := rand.Perm(max - min + 1)
	// fmt.Println(permutation)

	var subWiner []int
	for _, v := range permutation {
		if v != 0 {
			subWiner = append(subWiner, v)
		}
	}
	fmt.Println(subWiner)
	for i := 1; i <= int(totalWiner); i++ {
		fmt.Println(subWiner[i])
		database.DB.Where("raffle_id = ?", id).Where("queue_no = ?", subWiner[i]).Model(&participant).Updates(map[string]interface{}{"Status": 1, "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})
	}

	database.DB.Where("id = ?", id).Model(&raffle).Updates(map[string]interface{}{"Status": 1, "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})

	return c.JSON(fiber.Map{
		"total_participant": totalParticipant,
		"total_winer_max":   int(totalWiner),
		"status":            "200",
		"message":           "Success",
		"data":              "[]",
	})
}

func RaffleSearchParticipant(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var participant []models.Participant
	var data_raffle_id string
	var data_first_name string
	var data_last_name string
	var data_email string
	var data_instagram string
	var data_phone string
	var data_identity_no string
	var data_size string
	// tx := database.DB
	// raffle_id, nama, email,instagram, no telepon, id card no, size
	if data["raffle_id"] != "" {
		data_raffle_id = "participants.raffle_id =  " + data["raffle_id"] + ""
	} else {
		data_raffle_id = ""
	}
	if data["all_search"] != "" {
		data_first_name = "users.first_name LIKE  '%" + data["all_search"] + "%'"
	} else {
		data_first_name = ""
	}
	if data["all_search"] != "" {
		data_last_name = "users.last_name LIKE  '%" + data["all_search"] + "%'"
	} else {
		data_last_name = ""
	}
	if data["all_search"] != "" {
		data_email = "users.email LIKE '%" + data["all_search"] + "%'"
	} else {
		data_email = ""
	}
	if data["all_search"] != "" {
		data_instagram = "users.instagram LIKE '%" + data["all_search"] + "%'"
	} else {
		data_instagram = ""
	}
	if data["all_search"] != "" {
		data_phone = "users.phone LIKE '%" + data["all_search"] + "%'"
	} else {
		data_phone = ""
	}
	if data["all_search"] != "" {
		data_identity_no = "users.identity_no LIKE '%" + data["all_search"] + "%'"
	} else {
		data_identity_no = ""
	}
	if data["all_search"] != "" {
		data_size = "raffle_product_size_stocks.size LIKE '%" + data["all_search"] + "%'"
	} else {
		data_size = ""
	}

	tx := database.DB.Debug().
		Where(data_first_name).
		Or(data_last_name).
		Or(data_email).
		Or(data_instagram).
		Or(data_phone).
		Or(data_identity_no).
		Or(data_size).
		Where("participants.deleted_at IS NULL").
		Where(data_raffle_id).
		Joins("JOIN raffles on participants.raffle_id = raffles.id").
		Joins("JOIN users on participants.user_id = users.id").
		Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
		Joins("JOIN raffle_products on participants.raffle_product_id=raffle_products.id").
		Preload("Raffle").
		Preload("User").
		Preload("RaffleProduct").
		Preload("RaffleProductSizeStock").
		Find(&participant)

		// for _, ar := range participant {
		// 	fmt.Println(ar.User.FirstName)
		// }

		// fmt.Println(participant)
		// 	tx.Find(&participant)
	var total int64
	tx.Find(&participant).Count(&total)

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
		"data":    participant,
	})
}

func RaffleParticipantWinner(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id, _ := strconv.Atoi(data["raffle_id"])

	var participant []models.Participant

	// tx := database.DB

	tx := database.DB.Where("participants.raffle_id = ?", id).
		Where("participants.deleted_at IS NULL").
		Where("participants.status = 1").
		Joins("JOIN users on participants.user_id = users.id").
		Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
		Joins("JOIN raffle_products on participants.raffle_product_id=raffle_products.id").
		Preload("User").
		Preload("RaffleProduct").
		Preload("RaffleProductSizeStock").
		Find(&participant)
	// fmt.Println(tx)
	// for _, ar := range participant {
	// 	fmt.Println(ar.User.FirstName)
	// }

	// fmt.Println(participant)
	// 	tx.Find(&participant)
	var total int64
	tx.Find(&participant).Count(&total)

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
		"data":    participant,
	})
}

func RaffleParticipantUpdateStatus(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	raffle_id, _ := strconv.Atoi(data["raffle_id"])
	user_id, _ := strconv.Atoi(data["user_id"])
	status, _ := strconv.Atoi(data["status"])
	raffle_product_size_stock_id, _ := strconv.Atoi(data["raffle_product_size_stock_id"])

	//untuk dapetin stock barang
	var stock int
	database.DB.Raw("SELECT stock FROM raffle_product_size_stocks WHERE id = ?", raffle_product_size_stock_id).Scan(&stock)

	//untuk dapetin jumlah pemengan
	var winner_count int
	database.DB.Raw("SELECT COUNT(id) FROM participants WHERE status = 1 AND raffle_product_size_stock_id = ?", raffle_product_size_stock_id).Scan(&winner_count)
	if status == 1 {
		if stock <= winner_count {
			return c.JSON(fiber.Map{
				"status":  "202",
				"message": "Jumlah pemenang sudah cukup!",
				"data":    []string{},
			})
		}
	}
	var participant models.Participant
	database.DB.Where("raffle_id = ?", raffle_id).Where("user_id = ?", user_id).Model(&participant).Updates(map[string]interface{}{"Status": status, "UpdatedAt": time.Now(), "UpdatedBy": "Admin"})

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    []string{},
	})
}
func SetDatePay(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.Atoi(data["raffle_id"])
	var raffle models.Raffle
	database.DB.Where("id = ?", id).Model(&raffle).Updates(map[string]interface{}{"StartDatePay": data["start_date_pay"], "EndDatePay": data["end_date_pay"], "UpdatedAt": time.Now(), "UpdatedBy": "Start and End Date Pay"})
	const (
		DDMMYYYYhhmmss = "02-01-2006 15:04:05"
	)

	var participantUpdate models.Participant
	database.DB.Debug().Where("raffle_id = ?", id).Where("status = ?", 1).Where("status_transaction = ?", 0).Model(&participantUpdate).Updates(map[string]interface{}{"StartDatePay": data["start_date_pay"], "EndDatePay": data["end_date_pay"], "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Succes",
		"data":    "[]",
	})

}

func EmailBlastWinner(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.Atoi(data["raffle_id"])

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
		Where("participants.status = ?", 1).
		Where("participants.status_transaction = ?", 0).
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
			"title": "Pemenang Raffle",` + "" + `
			"body":  "<h2>Selamat ` + dataParticipant.User.FirstName + ` - ` + dataParticipant.User.LastName + ` Kamu terpilih menjadi pemenang Raffle ` + raffle.Name + `!</h2>` + "" +
			`<h3>Nama Produk : <a href='` + dataParticipant.RaffleProductSizeStock.UrlProduct + `'>` + dataParticipant.RaffleProduct.Name + `</a>` + "" +
			`<h3> Pembayaran Maksimal dari tanggal ` + start_date_pay + ` sampai tanggal ` + end_date_pay + ` </h3>` + "" +
			`<h3> Tinggal sedikit lagi. Yuk, selesaikan pembayaranmu!</h3>"` + "" + `
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
	database.DB.Debug().Where("id = ?", id).Model(&raffle).Updates(map[string]interface{}{"Status": 2, "UpdatedAt": time.Now(), "UpdatedBy": "Generate By System"})
	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Susses",
		"data":    "[]",
	})

	// fmt.Println("HTTP JSON POST URL:", link_service_email)

	// var jsonDataDatePay, err = json.Marshal(&data)
	// reqBodyData := strings.NewReader(string(jsonDataDatePay))
	// fmt.Println("reqBodyData:", reqBodyData)

	//

	// return c.JSON(fiber.Map{
	// 	"status":        "200",
	// 	"message":       "Succes",
	// 	"error_message": "",
	// 	"data":          "[]",
	// })

	// return fmt.Println(string(responseData))

	// } else {

	// 	return c.JSON(fiber.Map{
	// 		"status":        "200",
	// 		"message":       "Succes",
	// 		"error_message": "Email Sudah Dikirim",
	// 		"data":          "[]",
	// 	})
	// }

}

func SendEmailPerWinner(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.Atoi(data["raffle_id"])
	participants_id, _ := strconv.Atoi(data["participants_id"])

	fmt.Println("raffle_id", id)
	status := os.Getenv("STATUS")
	var link_service_email string
	if status == "dev" {
		link_service_email = os.Getenv("EMAIL_SERVICE_DEV")
	} else {
		link_service_email = os.Getenv("EMAIL_SERVICE_PROD")
	}

	var raffle models.Raffle
	database.DB.Where("id = ?", id).Find(&raffle)
	const (
		DDMMYYYYhhmmss = "02-01-2006 15:04:05"
	)
	ts_start_date_pay := raffle.StartDatePay

	start_date_pay := ts_start_date_pay.Format(DDMMYYYYhhmmss)
	fmt.Println("data end_date_pay :", start_date_pay)

	ts_end_date_pay := raffle.EndDatePay

	end_date_pay := ts_end_date_pay.Format(DDMMYYYYhhmmss)
	fmt.Println("data end_date_pay :", end_date_pay)

	var participant models.Participant
	database.DB.Where("participants.id = ?", participants_id).
		Where("participants.deleted_at IS NULL").
		Where("participants.raffle_id = ?", id).
		Where("participants.status = ?", 1).
		Joins("JOIN users on participants.user_id = users.id").
		Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
		Joins("JOIN raffle_products on participants.raffle_product_id=raffle_products.id").
		Preload("User").
		Preload("RaffleProduct").
		Preload("RaffleProductSizeStock").
		Find(&participant)

		// fmt.Println("data Participant :", participant)

	fmt.Println("data Participant :", participant.User.Email)
	fmt.Println("data Participant :", participant.RaffleProductSizeStock.UrlProduct)

	payload := strings.NewReader(`{` + "" + `
			"email": "` + participant.User.Email + `", 
			"host": "jdsport",` + "" + `
			"title": "Pemenang Raffle",` + "" + `
			"body":  "<h2>Selamat ` + participant.User.FirstName + ` - ` + participant.User.LastName + ` Kamu terpilih menjadi pemenang Raffle ` + raffle.Name + `!</h2>` + "" +
		`<h3>Nama Produk : <a href='` + participant.RaffleProductSizeStock.UrlProduct + `'>` + participant.RaffleProduct.Name + `</a>` + "" +
		`<h3> Pembayaran Maksimal dari tanggal ` + start_date_pay + ` sampai tanggal ` + end_date_pay + ` </h3>` + "" +
		`<h3> Tinggal sedikit lagi. Yuk, selesaikan pembayaranmu!</h3>"` + "" + `
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

	return c.JSON(fiber.Map{
		"status":        "200",
		"message":       "Succes",
		"error_message": "",
		"data":          "[]",
	})

}

func GetAllParticipant(c *fiber.Ctx) error {
	var result []models.User
	database.DB.Table("users").Where("user_type = ?", "customer").Where("deleted_at is NULL").Scan(&result)
	// row.Scan(&result)

	var total int64
	row2 := database.DB.Table("users").Select("COUNT(id)").Where("user_type = ?", "customer").Where("deleted_at is NULL").Row()
	row2.Scan(&total)

	if total == 0 {
		return c.JSON(fiber.Map{
			"status":  "202",
			"message": "Data Not Found",
			"data":    result,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    result,
	})
}

func UpdateStatusParticipant(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := strconv.Atoi(data["user_id"])
	reason, _ := data["reason"]
	status, _ := strconv.Atoi(data["status"])
	current_time := time.Now()

	// update status
	database.DB.Table("users").
		Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "reason": reason, "updated_at": current_time, "updated_by": "Admin"})

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "Success",
		"data":    []string{},
	})
}

func GetAllModuleSearchParticipant(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var participant []models.Participant
	var data_raffle_id string
	var data_first_name string
	var data_last_name string
	var data_email string
	var data_instagram string
	var data_phone string
	var data_identity_no string
	var data_size string
	var data_status string
	var data_registration_date string
	// tx := database.DB
	// raffle_id, nama, email,instagram, no telepon, id card no, size
	if data["raffle_id"] != "" {
		data_raffle_id = "participants.raffle_id =  " + data["raffle_id"] + ""
	} else {
		data_raffle_id = ""
	}
	if data["all_search"] != "" {
		data_first_name = "users.first_name LIKE  '%" + data["all_search"] + "%'"
	} else {
		data_first_name = ""
	}
	if data["all_search"] != "" {
		data_last_name = "users.last_name LIKE  '%" + data["all_search"] + "%'"
	} else {
		data_last_name = ""
	}
	if data["all_search"] != "" {
		data_email = "users.email LIKE '%" + data["all_search"] + "%'"
	} else {
		data_email = ""
	}
	if data["all_search"] != "" {
		data_instagram = "users.instagram LIKE '%" + data["all_search"] + "%'"
	} else {
		data_instagram = ""
	}
	if data["all_search"] != "" {
		data_phone = "users.phone LIKE '%" + data["all_search"] + "%'"
	} else {
		data_phone = ""
	}
	if data["all_search"] != "" {
		data_identity_no = "users.identity_no LIKE '%" + data["all_search"] + "%'"
	} else {
		data_identity_no = ""
	}
	if data["all_search"] != "" {
		data_size = "raffle_product_size_stocks.size LIKE '%" + data["all_search"] + "%'"
	} else {
		data_size = ""
	}
	if data["status"] != "" {
		data_status = "users.status = " + data["status"] + ""
	} else {
		data_status = ""
	}
	if data["registration_date"] != "" {
		data_registration_date = "users.created_at LIKE '" + data["registration_date"] + "%'"
	} else {
		data_registration_date = ""
	}

	tx := database.DB.Debug().
		Where(data_first_name).
		Or(data_last_name).
		Or(data_email).
		Or(data_instagram).
		Or(data_phone).
		Or(data_identity_no).
		Or(data_size).
		Where("participants.deleted_at IS NULL").
		Where(data_raffle_id).
		Where(data_registration_date).
		Where(data_status).
		Joins("JOIN users on participants.user_id = users.id").
		Joins("JOIN raffle_product_size_stocks on participants.raffle_product_size_stock_id=raffle_product_size_stocks.id").
		Joins("JOIN raffle_products on participants.raffle_product_id=raffle_products.id").
		Preload("User").
		Preload("RaffleProduct").
		Preload("RaffleProductSizeStock").
		Find(&participant)

		// for _, ar := range participant {
		// 	fmt.Println(ar.User.FirstName)
		// }

		// fmt.Println(participant)
		// 	tx.Find(&participant)
	var total int64
	tx.Find(&participant).Count(&total)

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
		"data":    participant,
	})
}
