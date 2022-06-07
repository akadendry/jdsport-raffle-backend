package jdsport_raffle_backend

import (
	"time"
)

type DataRaffle struct {
	Id                    uint      `json:"id;primaryKey"`
	Name                  string    `json:"name"`
	StartDateRegistration time.Time `json:"start_date_registration"`
	EndDateRegistration   time.Time `json:"end_date_registration"`
	AnnouncementDate      time.Time `json:"announcement_date"`
	StartDatePay          time.Time `json:"start_date_pay"`
	EndDatePay            time.Time `json:"end_date_pay"`
	Banner                string    `json:"banner"`
	BannerMobile          string    `json:"banner_mobile"`
	Copyright             string    `json:"copyright"`
	Slug                  string    `json:"slug"`
	SlugNo                uint      `json:"slug_no"`
	Status                uint      `json:"status"`
	IsEmail               uint      `json:"is_email"`
	CreatedAt             time.Time `json:"created_at"`
	CreatedBy             string    `json:"created_by"`
	UpdatedAt             time.Time `json:"updated_at"`
	UpdatedBy             string    `json:"updated_by"`
	DeletedAt             time.Time `json:"deleted_at"`
	DeletedBy             string    `json:"deleted_by"`
	// RaffleProduct         RaffleProduct `json:"product" gorm:"foreignKey:raffle_id"`
}

// type DataDetailProduct struct {
// 	Id                     uint                   `json:"id;primaryKey"`
// 	RaffleId               uint                   `json:"raffle_id"`
// 	Name                   string                 `json:"name"`
// 	Description            string                 `json:"description"`
// 	Image                  string                 `json:"image"`
// 	ImageMobile            string                 `json:"image_mobile"`
// 	RaffleProductSizeStock RaffleProductSizeStock `json:"detail" gorm:"foreignKey:raffle_product_id"`
// }

// type DataDetailSizeStock struct {
// 	Id              uint   `json:"id"`
// 	RaffleProductId uint   `json:"raffle_product_id"`
// 	Size            string `json:"size"`
// 	Stock           string `json:"stock"`
// 	UrlProduct      string `json:"url_product"`
// 	Sku             string `json:"sku"`
// }
