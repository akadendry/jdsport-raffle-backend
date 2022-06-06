package models

import (
	"time"

	"gorm.io/gorm"
)

type RaffleProduct struct {
	Id                     uint                     `json:"id;primaryKey"`
	RaffleId               uint                     `json:"raffle_id"`
	Name                   string                   `json:"name"`
	Description            string                   `json:"description" gorm:"type:text"`
	Image                  string                   `json:"image"`
	ImageMobile            string                   `json:"image_mobile"`
	CreatedAt              time.Time                `json:"created_at" gorm:"DEFAULT:null"`
	CreatedBy              string                   `json:"created_by" gorm:"DEFAULT:null"`
	UpdatedAt              time.Time                `json:"updated_at" gorm:"DEFAULT:null"`
	UpdatedBy              string                   `json:"updated_by" gorm:"DEFAULT:null"`
	DeletedAt              time.Time                `json:"deleted_at" gorm:"DEFAULT:null"`
	DeletedBy              string                   `json:"deleted_by" gorm:"DEFAULT:null"`
	RaffleProductSizeStock []RaffleProductSizeStock `json:"detail" gorm:"foreignKey:raffle_product_id"`
}

func (raffleProduct *RaffleProduct) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&RaffleProduct{}).Count(&total)

	return total
}
