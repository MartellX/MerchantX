package models

import "gorm.io/gorm"

type Repository interface {
	GetDB() *gorm.DB
	SetDB(gdb *gorm.DB)
	NewOffer(offerId uint64, sellerId uint64, name string, price int64, quantity int, available bool) (*Offer, error)
	Update(o *Offer) error
	UpdateColumns(o *Offer, name string, price int64, quantity int, available bool) error
	Delete(o *Offer)
	FindOffersByConditions(args map[string]interface{}) ([]Offer, error)
	FindOffer(offerId, sellerId uint64) (*Offer, error)
}
