package models

import (
	"time"
)

type RaffleAdd struct {
	Id                    string              `json:"id"`
	Name                  string              `json:"name"`
	StartDateRegistration time.Time           `json:"start_date_registration"`
	EndDateRegistration   time.Time           `json:"end_date_registration"`
	AnnouncementDate      time.Time           `json:"announcement_date"`
	StartDatePay          time.Time           `json:"start_date_pay"`
	EndDatePay            time.Time           `json:"end_date_pay"`
	Banner                string              `json:"banner"`
	BannerMobile          string              `json:"banner_mobile"`
	Copyright             string              `json:"copyright"`
	Slug                  string              `json:"slug"`
	Product               []GetRaffleProducts `json:"product"`
}

type GetRaffleProducts struct {
	Id          string                      `json:"id"`
	RaffleId    string                      `json:"raffle_id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Image       string                      `json:"image"`
	ImageMobile string                      `json:"image_mobile"`
	Detail      []GetRaffleProductSizeStock `json:"detail"`
}

type GetRaffleProductSizeStock struct {
	Id              string `json:"id"`
	RaffleProductId string `json:"raffle_product_id"`
	Size            string `json:"size"`
	Stock           string `json:"stock"`
	UrlProduct      string `json:"url_product"`
	Sku             string `json:"sku"`
}
