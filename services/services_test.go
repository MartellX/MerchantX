package services_test

import (
	"MartellX/avito-tech-task/models"
	"MartellX/avito-tech-task/models/mocks"
	"MartellX/avito-tech-task/services"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"testing"
	"time"
)

func startTestdataServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.File("/testdata1", "./testdata/testdata1.xlsx")
	e.File("/testdata2", "./testdata/testdata2.xlsx")
	e.File("/testdata3", "./testdata/testdata3.xlsx")
	e.File("/testdata4", "./testdata/testdata4.xlsx")
	e.File("/emptydata", "./testdata/emptydata.xlsx")
	e.File("/badfile", "./testdata/badfile")
	e.Logger.Fatal(e.Start(":1234"))
}

func TestService_StartUploadingTask(t *testing.T) {
	go startTestdataServer()
	time.Sleep(10 * time.Millisecond)

	mockCtrl := gomock.NewController(t)
	g := NewWithT(t)

	cases := []struct {
		description string
		url         string
		sellerId    uint
		expect      func(repo *mocks.MockRepository)
		result      func(task *services.Task)
	}{
		{
			description: "9 created",
			sellerId:    123,
			url:         "http://localhost:1234/testdata1",
			expect: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindOffer(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).MaxTimes(9)
				repo.EXPECT().NewOffer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			result: func(task *services.Task) {
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				g.Expect(task.Info.Created).Should(Equal(9))
			},
		},
		{
			description: "9 updated",
			url:         "http://localhost:1234/testdata1",
			sellerId:    123,
			expect: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindOffer(gomock.AssignableToTypeOf(uint64(1)), gomock.AssignableToTypeOf(uint64(1))).Return(&models.Offer{}, nil).MaxTimes(9)
				repo.EXPECT().UpdateColumns(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(9)
			},
			result: func(task *services.Task) {
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				g.Expect(task.Info.Updated).Should(Equal(9))
			},
		},
		{
			description: "8 deleted, 1 ignored (was not in DB)",
			url:         "http://localhost:1234/testdata3",
			sellerId:    1,
			expect: func(repo *mocks.MockRepository) {
				gomock.InOrder(
					repo.EXPECT().FindOffer(uint64(1), gomock.AssignableToTypeOf(uint64(1))).Return(nil, gorm.ErrRecordNotFound),
					repo.EXPECT().FindOffer(gomock.AssignableToTypeOf(uint64(1)), gomock.AssignableToTypeOf(uint64(1))).Return(&models.Offer{}, nil).MaxTimes(8),
				)

				repo.EXPECT().Delete(gomock.AssignableToTypeOf(&models.Offer{})).MaxTimes(8)
			},
			result: func(task *services.Task) {
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				g.Expect(task.Info.Deleted).Should(Equal(8))
			},
		},
		{
			description: "6 errors",
			url:         "http://localhost:1234/testdata2",
			sellerId:    123,
			expect: func(repo *mocks.MockRepository) {

			},
			result: func(task *services.Task) {
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				g.Expect(task.Info.Errors).Should(Equal(6))
			},
		},
		{
			description: "3 created, 3 updated, 3 deleted, 5 errors",
			url:         "http://localhost:1234/testdata4",
			sellerId:    123,
			expect: func(repo *mocks.MockRepository) {
				gomock.InOrder(
					repo.EXPECT().FindOffer(gomock.AssignableToTypeOf(uint64(1)), gomock.AssignableToTypeOf(uint64(1))).
						Return(nil, gorm.ErrRecordNotFound).Times(3),
					repo.EXPECT().FindOffer(gomock.AssignableToTypeOf(uint64(1)), gomock.AssignableToTypeOf(uint64(1))).
						Return(&models.Offer{}, nil).Times(6),
				)
				repo.EXPECT().NewOffer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).Times(3)
				repo.EXPECT().UpdateColumns(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(3)
				repo.EXPECT().Delete(gomock.AssignableToTypeOf(&models.Offer{})).Times(3)

			},
			result: func(task *services.Task) {
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).ShouldNot(And(Equal(404), Equal(400)))
				g.Expect(task.Info.Created).Should(Equal(3))
				g.Expect(task.Info.Updated).Should(Equal(3))
				g.Expect(task.Info.Deleted).Should(Equal(3))
				g.Expect(task.Info.Errors).Should(Equal(5))
			},
		},
		{
			description: "Empty table",
			url:         "http://localhost:1234/emptydata",
			sellerId:    0,
			expect: func(repo *mocks.MockRepository) {

			},
			result: func(task *services.Task) {
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).Should(Equal(400))
			},
		},

		{
			description: "Bad url",
			url:         "aff",
			sellerId:    0,
			expect: func(repo *mocks.MockRepository) {

			},
			result: func(task *services.Task) {
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).Should(Equal(400))
			},
		},
		{
			description: "Url not accessible",
			url:         "http://notaccesbleurl.test/",
			sellerId:    0,
			expect: func(repo *mocks.MockRepository) {

			},
			result: func(task *services.Task) {
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).Should(Equal(400))
			},
		},
		{
			description: "Bad endpoint",
			url:         "http://localhost:1234/bad",
			sellerId:    0,
			expect: func(repo *mocks.MockRepository) {

			},
			result: func(task *services.Task) {
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).Should(Equal(400))
			},
		},
		{
			description: "Bad file",
			url:         "http://localhost:1234/bad",
			sellerId:    0,
			expect: func(repo *mocks.MockRepository) {

			},
			result: func(task *services.Task) {
				for task.StatusCode == 201 || task.StatusCode == 102 {
					time.Sleep(5 * time.Millisecond)
				}
				g.Expect(task.StatusCode).Should(Equal(400))
			},
		},
	}

	for _, c := range cases {
		repo := mocks.NewMockRepository(mockCtrl)
		c.expect(repo)
		service := services.NewService(repo)
		task, err := service.StartUploadingTask(c.sellerId, c.url)
		g.Expect(err).ShouldNot(HaveOccurred())
		c.result(task)
	}
}
