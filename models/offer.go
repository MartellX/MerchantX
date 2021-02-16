package models

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Offer struct {
	OfferId   uint64    `gorm:"primaryKey;autoIncrement:false" json:"offer_id"`
	SellerId  uint64    `gorm:"primaryKey;autoIncrement:false" json:"seller_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Quantity  int       `json:"quantity"`
	Available bool      `json:"available"`
}

func (r *PostgresRepository) NewOffer(offerId uint64, sellerId uint64, name string, price int64, quantity int, available bool) (*Offer, error) {
	offer := &Offer{OfferId: offerId, SellerId: sellerId, Name: name, Price: price, Quantity: quantity, Available: available}
	res := r.GetDB().Session(&gorm.Session{Logger: silentLogger}).Create(offer)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return offer, nil
}

func (r *PostgresRepository) Update(o *Offer) error {
	res := r.GetDB().Save(o)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *PostgresRepository) UpdateColumns(o *Offer, name string, price int64, quantity int, available bool) error {
	o.Name = name
	o.Price = price
	o.Quantity = quantity
	o.Available = available
	return r.Update(o)
}

func (r *PostgresRepository) Delete(o *Offer) {
	r.GetDB().Delete(o)
}

var silentLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		LogLevel: logger.Silent, // Log level
	},
)

//
//type OfferParams struct {
//	OfferId *uint
//	SellerId *uint
//	Name *string
//}

func (r *PostgresRepository) FindOffersByConditions(args map[string]interface{}) ([]Offer, error) {

	tx := r.GetDB().Session(&gorm.Session{Logger: silentLogger})
	var offers []Offer
	//conditions := make([]string, 0, 3)
	//conditionArgs := make([]interface{}, 0, 3)
	condition := tx.Model(&offers)
	if offerId, ok := args["offer_id"]; ok {

		// Если неправильного типа, то просто не добавляем в запрос, другое решение - возвращать ошибку
		switch offerId.(type) {
		case uint, uint64, uint32:
			//conditions = append(conditions, "offer_id = ?")
			//conditionArgs = append(conditionArgs, offerId)
			condition = condition.Where("offer_id = ?", offerId)
		}

	}
	if sellerId, ok := args["seller_id"]; ok {
		switch sellerId.(type) {
		case uint, uint64, uint32:
			//conditions = append(conditions, "seller_id = ?")
			//conditionArgs = append(conditionArgs, sellerId)
			condition = condition.Where("seller_id = ?", sellerId)
		}

	}
	if name, ok := args["name"]; ok {
		switch name.(type) {
		case string:
			//conditions = append(conditions, "name ~ ?")
			//conditionArgs = append(conditionArgs, name)
			condition = condition.Where("name ILIKE ?", "%"+name.(string)+"%")
		}

	}

	//result := GetDB().Find(&offers, strings.Join(conditions, "AND"), conditionArgs)

	result := condition.Find(&offers)

	if result.Error != nil {
		return nil, errors.New("database exception")
	}
	return offers, nil
}

func (r *PostgresRepository) FindOffer(offerId, sellerId uint64) (*Offer, error) {

	tx := r.GetDB().Session(&gorm.Session{Logger: silentLogger})
	var offer Offer

	result := tx.Where("offer_id = ? AND seller_id = ?", offerId, sellerId).First(&offer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &offer, nil
}
