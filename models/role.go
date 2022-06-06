package models
import (
	"time"
)
type Role struct {
	Id          uint         `json:"id"`
	Name        string       `json:"name"`
	Access []Access `json:"permissions" gorm:"many2many:role_accesses"`
	CreatedAt              time.Time                `json:"created_at" gorm:"DEFAULT:null"`
	CreatedBy              string                   `json:"created_by" gorm:"DEFAULT:null"`
	UpdatedAt              time.Time                `json:"updated_at" gorm:"DEFAULT:null"`
	UpdatedBy              string                   `json:"updated_by" gorm:"DEFAULT:null"`
	DeletedAt              time.Time                `json:"deleted_at" gorm:"DEFAULT:null"`
	DeletedBy              string                   `json:"deleted_by" gorm:"DEFAULT:null"`
}
