package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Offer struct {
	OfferId   uint	  `gorm:"primaryKey;autoIncrement:false" json:"offer_id"`
	SellerId  uint    `gorm:"primaryKey;autoIncrement:false" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Available bool    `json:"available"`
}

func NewOffer(offerId uint, sellerId uint, name string, price float64, quantity int, available bool) (*Offer, error) {
	offer := &Offer{OfferId: offerId, SellerId: sellerId, Name: name, Price: price, Quantity: quantity, Available: available}
	res := GetDB().Session(&gorm.Session{Logger: silentLogger}).Create(offer)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound{
		return nil, res.Error
	}
	return offer, nil
}


func (o *Offer) Update() error {
	res := GetDB().Save(o)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (o *Offer) UpdateColumns(name string, price float64, quantity int, available bool) error {
	o.Name = name
	o.Price = price
	o.Quantity = quantity
	o.Available = available
	return o.Update()
}

func (o *Offer) Delete()  {
	GetDB().Delete(o)
}

var silentLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		LogLevel:      logger.Silent, // Log level
},
)

//
//type OfferParams struct {
//	OfferId *uint
//	SellerId *uint
//	Name *string
//}

func FindOffersByConditions(args map[string]interface{}) ([]Offer, error) {

	tx := GetDB().Session(&gorm.Session{Logger: silentLogger})
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
			condition = condition.Where("name ~ ?", name)
		}

	}


	//result := GetDB().Find(&offers, strings.Join(conditions, "AND"), conditionArgs)

	result := condition.Find(&offers)

	if result.Error != nil {
		return nil, result.Error
	}
	return offers, nil
}



func FindOffer(offerId, sellerId uint) (*Offer, error) {

	tx := GetDB().Session(&gorm.Session{Logger: silentLogger})
	var offer Offer

	result := tx.Where("offer_id = ? AND seller_id = ?", offerId, sellerId).First(&offer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &offer, nil
}
