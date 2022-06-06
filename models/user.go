package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id                uint      `json:"id"`
	ErajayaClubUserId uint      `json:"erajaya_club_user_id" gorm:"unique;DEFAULT:null"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name" gorm:"DEFAULT:null"`
	Phone             string    `json:"phone" gorm:"DEFAULT:null"`
	Email             string    `json:"email"`
	IdentityNo        string    `json:"identity_no" gorm:"DEFAULT:null"`
	Password          []byte    `json:"-"`
	Instagram         string    `json:"instagram" gorm:"DEFAULT:null"`
	UserType          string    `json:"user_type"`
	RoleId            uint      `json:"role_id" gorm:"DEFAULT:null"`
	Status            int       `json:"status" gorm:"DEFAULT:1"`
	Reason            string    `json:"reason"`
	Token             string    `json:"token" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at" gorm:"DEFAULT:null"`
	CreatedBy         string    `json:"created_by" gorm:"DEFAULT:null"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"DEFAULT:null"`
	UpdatedBy         string    `json:"updated_by" gorm:"DEFAULT:null"`
	DeletedAt         time.Time `json:"deleted_at" gorm:"DEFAULT:null"`
	DeletedBy         string    `json:"deleted_by" gorm:"DEFAULT:null"`
	Role              Role      `json:"role" gorm:"foreignKey:RoleId"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)

	return total
}

func (user *User) Take(db *gorm.DB, limit int, offset int) interface{} {
	var products []User

	db.Preload("Role").Offset(offset).Limit(limit).Find(&products)

	return products
}
