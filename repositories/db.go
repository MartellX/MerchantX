package repositories

import (
	"MartellX/avito-tech-task/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"strings"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func NewRepositoryFromEnvironments() (*PostgresRepository, error) {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("db_host")

	if username == "" || password == "" || dbName == "" || dbHost == "" {
		return nil, nil
	}
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	//fmt.Println(dbUri)

	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	fmt.Println("Connected to database")
	db := conn
	db.AutoMigrate(&models.Offer{})

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgresRepository) SetDB(gdb *gorm.DB) {
	r.db = gdb
}

func (r *PostgresRepository) NewOffer(offerId uint64, sellerId uint64, name string, price int64, quantity int, available bool) (*models.Offer, error) {
	offer := &models.Offer{OfferId: offerId, SellerId: sellerId, Name: name, Price: price, Quantity: quantity, Available: available}
	res := r.GetDB().Session(&gorm.Session{Logger: silentLogger}).Create(offer)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return offer, nil
}

func (r *PostgresRepository) Update(o *models.Offer) error {
	res := r.GetDB().Save(o)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *PostgresRepository) UpdateColumns(o *models.Offer, name string, price int64, quantity int, available bool) error {
	o.Name = name
	o.Price = price
	o.Quantity = quantity
	o.Available = available
	return r.Update(o)
}

func (r *PostgresRepository) Delete(o *models.Offer) {
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

func (r *PostgresRepository) FindOffersByConditions(args map[string]interface{}) ([]models.Offer, error) {

	//tx := r.GetDB().Session(&gorm.Session{Logger: silentLogger})
	db, err := r.GetDB().DB()
	if err != nil {
		return nil, err
	}
	var offers []models.Offer

	conditions := make([]string, 0, 3)
	conditionArgs := make([]interface{}, 0, 3)

	//condition := tx.Model(&offers)
	if offerId, ok := args["offer_id"]; ok {

		// Если неправильного типа, то просто не добавляем в запрос, другое решение - возвращать ошибку
		switch offerId.(type) {
		case uint, uint64, uint32:
			conditions = append(conditions, "offer_id = $"+strconv.Itoa(len(conditions)+1))
			conditionArgs = append(conditionArgs, offerId)
			//condition = condition.Where("offer_id = ?", offerId)
		}

	}
	if sellerId, ok := args["seller_id"]; ok {
		switch sellerId.(type) {
		case uint, uint64, uint32:
			conditions = append(conditions, "seller_id = $"+strconv.Itoa(len(conditions)+1))
			conditionArgs = append(conditionArgs, sellerId)
			//condition = condition.Where("seller_id = ?", sellerId)
		}

	}
	if name, ok := args["name"]; ok {
		switch name.(type) {
		case string:
			conditions = append(conditions, "name ILIKE $"+strconv.Itoa(len(conditions)+1))
			conditionArgs = append(conditionArgs, "%"+name.(string)+"%")
			//condition = condition.Where("name ILIKE ?", "%"+name.(string)+"%")
		}

	}

	rows, err := db.Query(
		"SELECT * FROM \"offers\" WHERE "+
			strings.Join(conditions, " AND "), conditionArgs...)

	//result := condition.Find(&offers)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var offer models.Offer
		err := rows.Scan(
			&offer.OfferId, &offer.SellerId, &offer.CreatedAt, &offer.UpdatedAt,
			&offer.Name, &offer.Price, &offer.Quantity, &offer.Available)
		if err != nil {
			log.Fatal(err)
		}
		offers = append(offers, offer)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return offers, nil
}

func (r *PostgresRepository) FindOffer(offerId, sellerId uint64) (*models.Offer, error) {

	tx := r.GetDB().Session(&gorm.Session{Logger: silentLogger})
	var offer models.Offer

	result := tx.Where("offer_id = ? AND seller_id = ?", offerId, sellerId).First(&offer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &offer, nil
}
