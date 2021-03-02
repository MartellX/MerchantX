package models

import (
	"time"
)

type Offer struct {
	OfferId   uint64    `gorm:"primaryKey;autoIncrement:false;index:idx_offer" json:"offer_id"`
	SellerId  uint64    `gorm:"primaryKey;autoIncrement:false;index:idx_offer" json:"seller_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Quantity  int       `json:"quantity"`
	Available bool      `json:"available"`
}
