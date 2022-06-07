package jdsport_raffle_backend

import (
	"os"

	"github.com/akadendry/jdsport-raffle-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	// database, err := gorm.Open(mysql.Open("developer:dendry@/raffle"), &gorm.Config{})
	database, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       os.Getenv("DB_CONNECTION"), // data source name
		DefaultStringSize:         256,                        // default size for string fields
		DisableDatetimePrecision:  true,                       // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                       // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                       // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                      // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	DB = database

	database.AutoMigrate(&models.User{}, &models.Role{}, &models.Access{}, &models.RoleAccess{}, &models.RaffleProduct{}, &models.Raffle{}, &models.RaffleProductSizeStock{}, &models.Participant{}, &models.MessageServiceEmail{})
}
