package repository

import (
	"MartellX/avito-tech-task/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetDB() *gorm.DB
	SetDB(gdb *gorm.DB)
	NewOffer(offerId uint64, sellerId uint64, name string, price int64, quantity int, available bool) (*models.Offer, error)
	Update(o *models.Offer) error
	UpdateColumns(o *models.Offer, name string, price int64, quantity int, available bool) error
	Delete(o *models.Offer)
	FindOffersByConditions(args map[string]interface{}) ([]models.Offer, error)
	FindOffer(offerId, sellerId uint64) (*models.Offer, error)
}
