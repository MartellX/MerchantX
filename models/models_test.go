package models_test

import (
	"MartellX/avito-tech-task/models"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)



func SetNewMock() (sqlmock.Sqlmock, error){

	var mock sqlmock.Sqlmock

	var db *sql.DB
	var err error

	db, mock, err = sqlmock.New() // mock sql.DB
	if err != nil {
		return nil, err
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db}),
		&gorm.Config{}) // open gorm db
	if err != nil {
		return nil, err
	}
	models.SetDB(gdb)

	return mock, nil
}

func TestNewOffer(t *testing.T) {
	// Before
	g := NewGomegaWithT(t)
	mock, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffer := models.Offer{
		OfferId:   1,
		SellerId:  1,
		Name:      "abc",
		Price:     123.4,
		Quantity:  5,
		Available: true,
	}

	mock.
		ExpectExec(regexp.QuoteMeta("UPDATE \"offers\"")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	offer, err := models.NewOffer(testOffer.OfferId, testOffer.SellerId, testOffer.Name, testOffer.Price, testOffer.Quantity, testOffer.Available)

	g.Expect(err).ShouldNot(HaveOccurred())
	testOffer.CreatedAt = offer.CreatedAt
	testOffer.UpdatedAt = offer.UpdatedAt
	g.Expect(*offer).Should(Equal(testOffer))

	// After
	err = mock.ExpectationsWereMet()
	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestFindOffer(t *testing.T) {

	// Before
	g := NewGomegaWithT(t)
	mock, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffers := []models.Offer {
		{
			OfferId:   1,
			SellerId:  1,
			Name:      "abc",
			Price:     123.4,
			Quantity:  5,
			Available: true,
		},
		{
			OfferId:   3,
			SellerId:  2,
			Name:      "",
			Price:     0,
			Quantity:  12,
			Available: true,
		},
		{
			OfferId:   4542,
			SellerId:  11,
			Name:      "s",
			Price:     100,
			Quantity:  1,
			Available: true,
		},
	}

	for _, testOffer := range testOffers {
		rows := mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"}).
			AddRow(testOffer.OfferId, testOffer.SellerId, testOffer.CreatedAt, testOffer.UpdatedAt, testOffer.Name, testOffer.Price, testOffer.Quantity, testOffer.Available)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\"")).
			WithArgs(testOffer.OfferId, testOffer.SellerId).
			WillReturnRows(rows)

		offer, err := models.FindOffer(testOffer.OfferId, testOffer.SellerId)
		g.Expect(err).ShouldNot(HaveOccurred())
		g.Expect(*offer).Should(Equal(testOffer))
	}


	// After
	err = mock.ExpectationsWereMet()
	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestFindOffersByConditions(t *testing.T) {
	// Before
	g := NewGomegaWithT(t)
	mock, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffers := []models.Offer {
		{
			OfferId:   1,
			SellerId:  1,
			Name:      "iPhone 12",
			Price:     60000,
			Quantity:  50,
			Available: true,
		},
		{
			OfferId:   2,
			SellerId:  1,
			Name:      "Xiaomi",
			Price:     100,
			Quantity:  1,
			Available: true,
		},
		{
			OfferId:   3,
			SellerId:  1,
			Name:      "PC",
			Price:     40000,
			Quantity:  5,
			Available: true,
		},
		{
			OfferId:   3,
			SellerId:  2,
			Name:      "iPhone 11",
			Price:     50000,
			Quantity:  12,
			Available: true,
		},
		{
			OfferId:   4542,
			SellerId:  2,
			Name:      "Guitar",
			Price:     100,
			Quantity:  1,
			Available: true,
		},
		{
			OfferId:   444,
			SellerId:  5,
			Name:      "High end PC",
			Price:     200000,
			Quantity:  1,
			Available: true,
		},
	}

	testSqlRows := make([]*sqlmock.Rows, 0, len(testOffers))
	for _, offer := range testOffers {
		rows := mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"}).
			AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
		testSqlRows = append(testSqlRows, rows)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE seller_id = $1")).
		WithArgs(1).
		WillReturnRows(testSqlRows[:3]...)

	args := map[string]interface{}{
		"seller_id": uint(1),
	}
	offers, err := models.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(3), ContainElements(testOffers[:3])))


	// After
	err = mock.ExpectationsWereMet()
	g.Expect(err).ShouldNot(HaveOccurred())
}

