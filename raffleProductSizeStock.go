package jdsport_raffle_backend

import (
	"time"

	"gorm.io/gorm"
)

type RaffleProductSizeStock struct {
	Id              uint      `json:"id;primaryKey"`
	RaffleProductId uint      `json:"raffle_product_id"`
	Size            string    `json:"size"`
	Stock           string    `json:"stock"`
	UrlProduct      string    `json:"url_product"`
	Sku             string    `json:"sku"`
	CreatedAt       time.Time `json:"created_at" gorm:"DEFAULT:null"`
	CreatedBy       string    `json:"created_by" gorm:"DEFAULT:null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"DEFAULT:null"`
	UpdatedBy       string    `json:"updated_by" gorm:"DEFAULT:null"`
	DeletedAt       time.Time `json:"deleted_at" gorm:"DEFAULT:null"`
	DeletedBy       string    `json:"deleted_by" gorm:"DEFAULT:null"`
}

func (raffleProductSizeStock *RaffleProductSizeStock) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&RaffleProductSizeStock{}).Count(&total)

	return total
}
