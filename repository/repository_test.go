package repository_test

import (
	"MartellX/avito-tech-task/models"
	"MartellX/avito-tech-task/repository"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func SetNewMock() (sqlmock.Sqlmock, *repository.PostgresRepository, error) {

	var mock sqlmock.Sqlmock

	var db *sql.DB
	var err error

	db, mock, err = sqlmock.New() // mock sql.DB
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db}),
		&gorm.Config{}) // open gorm db
	if err != nil {
		return nil, nil, err
	}
	repository.NewRepository(gdb)

	return mock, repository.NewRepository(gdb), nil
}

func TestNewOffer(t *testing.T) {
	// Before
	g := NewGomegaWithT(t)
	mock, repo, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffer := models.Offer{
		OfferId:   1,
		SellerId:  1,
		Name:      "abc",
		Price:     1234,
		Quantity:  5,
		Available: true,
	}

	mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO \"offers\"")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	offer, err := repo.NewOffer(testOffer.OfferId, testOffer.SellerId, testOffer.Name, testOffer.Price, testOffer.Quantity, testOffer.Available)

	g.Expect(err).ShouldNot(HaveOccurred())
	testOffer.CreatedAt = offer.CreatedAt
	testOffer.UpdatedAt = offer.UpdatedAt
	g.Expect(*offer).Should(Equal(testOffer))

	// After
	err = mock.ExpectationsWereMet()
	g.Expect(err).ShouldNot(HaveOccurred())
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestOffer_UpdateColumns(t *testing.T) {
	// Before
	g := NewGomegaWithT(t)
	mock, repo, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffer := models.Offer{
		OfferId:   1,
		SellerId:  1,
		Name:      "abc",
		Price:     1234,
		Quantity:  5,
		Available: true,
	}

	updatingOffer := models.Offer{
		OfferId:   1,
		SellerId:  1,
		Name:      "yo",
		Price:     4321,
		Quantity:  1,
		Available: false,
	}

	mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO \"offers\"")).
		WithArgs(testOffer.OfferId, testOffer.SellerId, AnyTime{}, AnyTime{}, testOffer.Name, testOffer.Price, testOffer.Quantity, testOffer.Available).
		WillReturnResult(sqlmock.NewResult(int64(testOffer.OfferId), 1))

	offer, err := repo.NewOffer(testOffer.OfferId, testOffer.SellerId, testOffer.Name, testOffer.Price, testOffer.Quantity, testOffer.Available)
	g.Expect(err).ShouldNot(HaveOccurred())

	mock.
		ExpectExec(regexp.QuoteMeta("UPDATE \"offers\"")).
		WithArgs(AnyTime{}, AnyTime{}, updatingOffer.Name, updatingOffer.Price, updatingOffer.Quantity, updatingOffer.Available, offer.OfferId, offer.SellerId).
		WillReturnResult(sqlmock.NewResult(int64(offer.OfferId), 1))

	err = repo.UpdateColumns(offer, updatingOffer.Name, updatingOffer.Price, updatingOffer.Quantity, updatingOffer.Available)
	g.Expect(err).ShouldNot(HaveOccurred())

	// After
	err = mock.ExpectationsWereMet()
	g.Expect(err).ShouldNot(HaveOccurred())

}

func TestFindOffer(t *testing.T) {

	// Before
	g := NewGomegaWithT(t)
	mock, repo, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffers := []models.Offer{
		{
			OfferId:   1,
			SellerId:  1,
			Name:      "abc",
			Price:     1234,
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

		offer, err := repo.FindOffer(testOffer.OfferId, testOffer.SellerId)
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
	mock, repo, err := SetNewMock()
	g.Expect(err).ShouldNot(HaveOccurred())

	// Test
	testOffers := []models.Offer{
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

	// When seller_id is provided only him offers should return

	// seller_id = 1
	rows := mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range testOffers[:3] {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE seller_id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	args := map[string]interface{}{
		"seller_id": uint(1),
	}
	offers, err := repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(3), ContainElements(testOffers[:3])))

	// When offer_id is provided only offers with this id should return

	// offer_id = 3
	rows = mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range testOffers[3:5] {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE offer_id = $1")).
		WithArgs(3).
		WillReturnRows(rows)

	args = map[string]interface{}{
		"offer_id": uint(3),
	}
	offers, err = repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(2), ContainElements(testOffers[3:5])))

	// When name is provided only offers, which contains it, should return

	// name = iPhone

	rows = mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range append(testOffers[:1], testOffers[3]) {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE name ILIKE $1")).
		WithArgs("%iPhone%").
		WillReturnRows(rows)

	args = map[string]interface{}{
		"name": "iPhone",
	}
	offers, err = repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(2), ContainElements(append(testOffers[:1], testOffers[3]))))

	// Combining offer_id and seller_id

	// offer_id = 3 AND seller_id = 2

	rows = mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range testOffers[3:4] {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE offer_id = $1 AND seller_id = $2")).
		WithArgs(3, 2).
		WillReturnRows(rows)

	args = map[string]interface{}{
		"offer_id":  uint(3),
		"seller_id": uint(2),
	}
	offers, err = repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(1), ContainElements(testOffers[3:4])))

	// Combining offer_id and name

	// offer_id = 1 AND seller_id = iPhone

	rows = mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range testOffers[:1] {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE offer_id = $1 AND name ILIKE $2")).
		WithArgs(1, "%iPhone%").
		WillReturnRows(rows)

	args = map[string]interface{}{
		"offer_id": uint(1),
		"name":     "iPhone",
	}
	offers, err = repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(1), ContainElements(testOffers[:1])))

	// Combining seller_id and name

	// seller_id = 2 AND name = "Guit"

	rows = mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range testOffers[4:5] {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE seller_id = $1 AND name ILIKE $2")).
		WithArgs(2, "%Guit%").
		WillReturnRows(rows)

	args = map[string]interface{}{
		"seller_id": uint(2),
		"name":      "Guit",
	}
	offers, err = repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(1), ContainElements(testOffers[4:5])))

	// Combining all params

	// offer_id = 444 AND seller_id = 5 AND name = "PC"

	rows = mock.NewRows([]string{"offer_id", "seller_id", "created_at", "updated_at", "name", "price", "quantity", "available"})
	for _, offer := range testOffers[5:] {
		rows.AddRow(offer.OfferId, offer.SellerId, offer.CreatedAt, offer.UpdatedAt, offer.Name, offer.Price, offer.Quantity, offer.Available)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"offers\" WHERE offer_id = $1 AND seller_id = $2 AND name ILIKE $3")).
		WithArgs(444, 5, "%PC%").
		WillReturnRows(rows)

	args = map[string]interface{}{
		"offer_id":  uint(444),
		"seller_id": uint(5),
		"name":      "PC",
	}
	offers, err = repo.FindOffersByConditions(args)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(offers).Should(And(HaveLen(1), ContainElements(testOffers[5:])))

	// After
	err = mock.ExpectationsWereMet()
	g.Expect(err).ShouldNot(HaveOccurred())
}
