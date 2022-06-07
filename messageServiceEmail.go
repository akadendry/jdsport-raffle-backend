package jdsport_raffle_backend

import (
	"time"
)

type MessageServiceEmail struct {
	Id           uint      `json:"id;primaryKey"`
	TypeEmail    string    `json:"type_email"`
	LinkMans     string    `json:"link_mans"`
	LinkWomens   string    `json:"link_womens"`
	LinkKids     string    `json:"link_kids"`
	LinkBrands   string    `json:"link_brands"`
	MessageEmail string    `json:"message_email" gorm:"type:longtext"`
	CreatedAt    time.Time `json:"created_at" gorm:"DEFAULT:null"`
	CreatedBy    string    `json:"created_by" gorm:"DEFAULT:null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"DEFAULT:null"`
	UpdatedBy    string    `json:"updated_by" gorm:"DEFAULT:null"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"DEFAULT:null"`
	DeletedBy    string    `json:"deleted_by" gorm:"DEFAULT:null"`
}
