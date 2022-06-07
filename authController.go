package jdsport_raffle_backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akadendry/jdsport-raffle-backend/database"
	"github.com/akadendry/jdsport-raffle-backend/models"
	"github.com/akadendry/jdsport-raffle-backend/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email Or Password do not match",
		})
	}
	role_id, _ := strconv.Atoi(data["role_id"])

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    uint(role_id),
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).Where("user_type = ?", "admin").Where("deleted_at IS NULL").First(&user)

	if user.Id == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Email Or Password not found",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		// c.Status(400)
		return c.JSON(fiber.Map{
			"status":  "400",
			"message": "Incorrect Email Or Password",
		})
	}

	if user.Id != 0 && user.Status == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "202",
			"message": "User Inactive, Please contact Administrator",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"data":    user,
		"status":  "200",
		"message": "success",
	})
}

func LoginMember(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	httpposturl_login := os.Getenv("HTTPPOSTURL_LOGIN")
	httpposturl_profile := os.Getenv("HTTPPOSTURL_PROFILE")
	// fmt.Println("HTTP JSON POST URL:", httpposturl_login)

	var jsonDataLogin, err = json.Marshal(&data)
	reqBodyData := strings.NewReader(string(jsonDataLogin))
	response_login, err := http.NewRequest("POST", httpposturl_login, reqBodyData)
	response_login.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response_login.Header.Add("Device-Id", "63e5f33e-d763-43ed-9cb4-57d933dfaf50")
	response_login.Header.Add("Source", "eraspace")

	client := &http.Client{}
	responsed, error := client.Do(response_login)
	if error != nil {
		panic(error)
	}
	defer responsed.Body.Close()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(responsed.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("test jose 36 : ",responseData)
	var responseObject models.ResponseLogin
	json.Unmarshal(responseData, &responseObject)

	loginToken := responseObject.Datas.Token

	if loginToken == "" {
		return c.JSON(fiber.Map{
			"status":        "404",
			"dataSSO":       []string{},
			"dataProfilSSO": []string{},
			"data":          []string{},
			"message":       "Username dan password salah!",
		})
	}

	response_profil, err := http.NewRequest("GET", httpposturl_profile, strings.NewReader(``))
	response_profil.Header.Add("Device-Id", "df9e0c5f-e267-427d-a890-31c10641fec1")
	response_profil.Header.Add("Source", "eraspace")
	response_profil.Header.Add("Content-Type", "application/json")
	response_profil.Header.Add("Device", "Postman")
	response_profil.Header.Add("Authorization", "Bearer "+loginToken)
	// fmt.Println("Bearer " + loginToken)
	client_profil := &http.Client{}
	responsed_profil, error := client_profil.Do(response_profil)

	if error != nil {
		panic(error)
	}
	defer responsed_profil.Body.Close()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseDataProfil, err := ioutil.ReadAll(responsed_profil.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObjectProfil models.ResponseProfil
	json.Unmarshal(responseDataProfil, &responseObjectProfil)

	user := models.User{
		ErajayaClubUserId: responseObjectProfil.Datas.ID,
		FirstName:         responseObjectProfil.Datas.Firstname,
		LastName:          responseObjectProfil.Datas.Lastname,
		Phone:             responseObjectProfil.Datas.Phone,
		Email:             responseObjectProfil.Datas.Email,
		IdentityNo:        responseObjectProfil.Datas.IdentityNumber,
		Instagram:         responseObjectProfil.Datas.Attributes.SocialMedia.Instagram,
		UserType:          "customer",
		CreatedAt:         time.Now(),
		CreatedBy:         "Create SSO",
	}

	var participant []models.Participant
	raffle_id, _ := strconv.Atoi(data["raffle_id"])
	var total int64
	database.DB.Find(&participant, "erajaya_club_user_id = ? AND raffle_id = ? ", responseObjectProfil.Datas.ID, raffle_id).Count(&total)

	if total >= 1 {
		return c.JSON(fiber.Map{
			"status":        "202",
			"dataSSO":       []string{},
			"dataProfilSSO": []string{},
			"data":          []string{},
			"message":       "User sudah pernah terdaftar pada raffle ini!",
		})
	}

	//sampe sini cobain array push string
	var data_update []string
	if responseObjectProfil.Datas.Firstname != "" {
		data_update = append(data_update, "first_name")
	}
	if responseObjectProfil.Datas.Lastname != "" {
		data_update = append(data_update, "last_name")
	}
	if responseObjectProfil.Datas.Phone != "" {
		data_update = append(data_update, "phone")
	}
	if responseObjectProfil.Datas.Email != "" {
		data_update = append(data_update, "email")
	}
	if responseObjectProfil.Datas.IdentityNumber != "" {
		data_update = append(data_update, "identity_no")
	}
	if responseObjectProfil.Datas.Attributes.SocialMedia.Instagram != "" {
		data_update = append(data_update, "instagram")
	}

	//titik akhir coba

	database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "erajaya_club_user_id"}}, // key colume
		DoUpdates: clause.AssignmentColumns(data_update),           // column needed to be updated
	}).Create(&user)

	var dataUser models.User
	database.DB.Where("erajaya_club_user_id = ?", responseObjectProfil.Datas.ID).Find(&dataUser)

	// userId, _ := strconv.Atoi(dataUser.Id)
	return c.JSON(fiber.Map{
		"status":        "200",
		"dataSSO":       responseObject,
		"dataProfilSSO": responseObjectProfil,
		"data":          dataUser,
		"message":       "success",
	})

}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)
	database.DB.Preload("Role").Find(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email Or Password do not match",
		})
	}

	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id: uint(userId),
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func ResetPasswordAdmin(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).Where("user_type = ?", "admin").First(&user)

	if user.Id == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Email not found",
		})
	}

	if user.Id != 0 && user.Status == 0 {
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "202",
			"message": "User Inactive, Please contact Administrator",
		})
	}

	if user.Id >= 0 && user.Status == 1 {
		StringRandom := StringRandom(64)

		userDb := models.User{
			Id:    uint(user.Id),
			Token: string(StringRandom),
		}
		database.DB.Model(&userDb).Updates(userDb)

		to := []string{data["email"]}
		cc := []string{""}
		subject := "Reset Password"
		message := "Hey <b>" + user.FirstName + " " + user.LastName + "</b> <br>" +
			"Permohonan reset password Anda telah kami terima " +
			"<br>" +
			"<br>" +
			"<a href='" + os.Getenv("APP_URL") + "/cms/account/password/reset?token=" + StringRandom + "'>RESET</a>" + "<br>" +
			"<br>" +
			"Email dibuat secara otomatis. Mohon tidak mengirimkan balasan ke email ini."
		attach := ""

		err := sendMail(to, cc, subject, message, attach)
		if err != nil {
			log.Fatal(err.Error())
		}

		// fmt.Println("test : " + StringRandom(8))
		// c.Status(404)
		return c.JSON(fiber.Map{
			"status":  "200",
			"message": "Success Reset Password",
		})
	}

	// log.Println("Mail sent!")
	return c.JSON("Mail sent!")
}

func GetUserByToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	var total int64
	database.DB.Find(&user, "token = ?", data["token"]).Count(&total)

	if total == 0 || data["token"] == "" {
		return c.JSON(fiber.Map{
			"status":  "404",
			"data":    "[]",
			"message": "Token verifikasi tidak valid, silakan cek ulang email Anda",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "200",
		"data":    user,
		"message": "success",
	})

}

func UpdatePasswordAdmin(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	var getUser models.User
	database.DB.Where("email = ?", data["email"]).Where("user_type = ?", "admin").First(&getUser)

	if data["email"] == os.DevNull || data["email"] == "" {
		return c.JSON(fiber.Map{
			"status":  "404",
			"message": "Fail Update Password",
		})
	}
	user := models.User{
		Id:    uint(getUser.Id),
		Token: "",
	}

	database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},            // key colume
		DoUpdates: clause.AssignmentColumns([]string{"token"}), // column needed to be updated
	}).Create(&user)

	user.SetPassword(data["password"])
	database.DB.Model(&user).Updates(user)

	return c.JSON(fiber.Map{
		"status":  "200",
		"message": "success",
	})
}
