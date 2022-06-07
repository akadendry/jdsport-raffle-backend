package jdsport_raffle_backend

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/akadendry/jdsport-raffle-backend/v2/database"
	"github.com/akadendry/jdsport-raffle-backend/v2/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func AllCustomerUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB.Where("deleted_at IS NULL").Where("user_type = ?", "customer"), &models.User{}, page))
}

func AllAdminUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(database.DB.Where("deleted_at IS NULL").Where("user_type = ?", "admin"), &models.User{}, page))
}

func CreateUserNew(c *fiber.Ctx) error {
	// payloads, _ := ioutil.ReadAll(r.Body)

	// var user models.User
	// json.Unmarshal(payloads, &user)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	idErajayaClubUser, _ := strconv.Atoi(data["erajaya_club_user_id"])
	roleId, _ := strconv.Atoi(data["role_id"])
	user := models.User{
		ErajayaClubUserId: uint(idErajayaClubUser),
		FirstName:         data["first_name"],
		LastName:          data["last_name"],
		Phone:             data["phone"],
		Email:             data["email"],
		IdentityNo:        data["identity_no"],
		Instagram:         data["instagram"],
		UserType:          "customer",
		RoleId:            uint(roleId),
		CreatedAt:         time.Now(),
		CreatedBy:         data["erajaya_club_user_id"],
	}

	database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "erajaya_club_user_id"}},                                                             // key colume
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "phone", "email", "identity_no", "instagram"}), // column needed to be updated
	}).Create(&user)

	// database.DB.Create(&user)

	res := Result{Code: 200, Data: user, Message: "Success"}
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

func CreateUserAdmin(c *fiber.Ctx) error {

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

	user := models.User{
		FirstName:  data["first_name"],
		LastName:   data["last_name"],
		Phone:      data["phone"],
		Email:      data["email"],
		IdentityNo: data["identity_no"],
		Instagram:  data["instagram"],
		UserType:   "admin",
		CreatedAt:  time.Now(),
		CreatedBy:  data["erajaya_club_user_id"],
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	database.DB.Preload("Role").Find(&user)

	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id:        uint(id),
		DeletedAt: time.Now(),
	}

	database.DB.Model(&user).Updates(user)

	return nil
}
