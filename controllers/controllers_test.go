package controllers_test

import (
	"MartellX/avito-tech-task/controllers"
	"MartellX/avito-tech-task/models"
	"MartellX/avito-tech-task/repositories/mock_repositories"
	"MartellX/avito-tech-task/services"
	"MartellX/avito-tech-task/services/mock_services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_NewTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	g := NewWithT(t)
	e := echo.New()

	cases := []struct {
		description string
		expect      func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository)
	}{
		{
			description: "if all params provided - returning new task",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("seller_id", "1")
				f.Set("url", "https://example.com")
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				task := services.Task{StatusCode: http.StatusCreated}
				s.EXPECT().StartUploadingTask(gomock.AssignableToTypeOf(uint64(1)), gomock.AssignableToTypeOf("str")).
					Return(&task, nil)

				h := controllers.NewHandler(s, r)

				g.Expect(h.NewTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusCreated))
				var returnedTask services.Task
				g.Expect(json.Unmarshal(rec.Body.Bytes(), &returnedTask)).ShouldNot(HaveOccurred())
				g.Expect(returnedTask).Should(Equal(task))
			},
		},
		{
			description: "If seller_id not provided - return code 400 and inform about that parameter",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("url", "https://example.com")
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				rec := httptest.NewRecorder()
				h := controllers.NewHandler(s, r)
				c := e.NewContext(req, rec)

				g.Expect(h.NewTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusBadRequest))
				g.Expect(rec.Body).Should(ContainSubstring("seller_id"))
			},
		},

		{
			description: "If url not provided - return code 400 and inform about that parameter",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("seller_id", "123")
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				rec := httptest.NewRecorder()
				h := controllers.NewHandler(s, r)
				c := e.NewContext(req, rec)

				g.Expect(h.NewTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusBadRequest))
				g.Expect(rec.Body).Should(ContainSubstring("url"))
			},
		},

		{
			description: "If seller_id is not number - return 400 and error message",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("seller_id", "asd")
				f.Set("url", "https://example.com")
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				rec := httptest.NewRecorder()
				h := controllers.NewHandler(s, r)
				c := e.NewContext(req, rec)

				g.Expect(h.NewTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusBadRequest))
				g.Expect(rec.Body.String()).Should(ContainSubstring("message"))
			},
		},
	}

	for _, c := range cases {
		s := mock_services.NewMockTaskService(mockCtrl)
		r := mock_repositories.NewMockRepository(mockCtrl)
		fmt.Println(c.description)
		c.expect(s, r)
		fmt.Println("ok")

	}
}

func TestHandler_GetTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	g := NewWithT(t)
	e := echo.New()

	cases := []struct {
		description string
		expect      func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository)
	}{
		{
			description: "returning task by given id",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("task_id", "1")
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				task := services.Task{}
				s.EXPECT().GetTask(gomock.AssignableToTypeOf("string")).Return(&task, true)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusOK))
				var returnedTask services.Task
				g.Expect(json.Unmarshal(rec.Body.Bytes(), &returnedTask)).ShouldNot(HaveOccurred())
				g.Expect(returnedTask).Should(Equal(task))
			},
		},
		{
			description: "if task is not found - return error with NotFound code",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("task_id", "1")
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				s.EXPECT().GetTask(gomock.AssignableToTypeOf("string")).Return(nil, false)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusNotFound))
			},
		},
		{
			description: "if task_id is not provided - return error and inform",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetTask(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusBadRequest))
				g.Expect(rec.Body.String()).Should(ContainSubstring("task_id"))
			},
		},
	}

	for _, c := range cases {
		s := mock_services.NewMockTaskService(mockCtrl)
		r := mock_repositories.NewMockRepository(mockCtrl)
		fmt.Println(c.description)
		c.expect(s, r)
		fmt.Println("ok")

	}
}

func TestHandler_GetOffers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	g := NewWithT(t)
	e := echo.New()

	cases := []struct {
		description string
		expect      func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository)
	}{
		{
			description: "If seller_id, offer_id, name provided -> returning response with 3 offers",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("offer_id", "1")
				f.Set("seller_id", "1")
				f.Set("name", "example")
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				expectingArgs := map[string]interface{}{
					"offer_id":  uint64(1),
					"seller_id": uint64(1),
					"name":      "example",
				}
				r.EXPECT().FindOffersByConditions(expectingArgs).Return(make([]models.Offer, 3), nil)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetOffers(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusOK))
				g.Expect(gjson.GetBytes(rec.Body.Bytes(), "count").Int()).Should(BeEquivalentTo(3))
				g.Expect(gjson.GetBytes(rec.Body.Bytes(), "items").Array()).Should(HaveLen(3))

			},
		},
		{
			description: "If nothing provided -> returning response with 3 offers",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				expectingArgs := map[string]interface{}{}
				r.EXPECT().FindOffersByConditions(expectingArgs).Return(make([]models.Offer, 3), nil)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetOffers(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusOK))
				g.Expect(gjson.GetBytes(rec.Body.Bytes(), "count").Int()).Should(BeEquivalentTo(3))
				g.Expect(gjson.GetBytes(rec.Body.Bytes(), "items").Array()).Should(HaveLen(3))

			},
		},
		{
			description: "If seller_id is wrong format -> return error message",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("seller_id", "s")
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetOffers(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusBadRequest))
				g.Expect(gjson.GetBytes(rec.Body.Bytes(), "message").Str).Should(ContainSubstring("seller_id"))
			},
		},
		{
			description: "If offer_id is wrong format -> return error message",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				f.Set("offer_id", "-2")
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetOffers(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusBadRequest))
				g.Expect(gjson.GetBytes(rec.Body.Bytes(), "message").Str).Should(ContainSubstring("offer_id"))
			},
		},
		{
			description: "If error occurred -> return 500 code",
			expect: func(s *mock_services.MockTaskService, r *mock_repositories.MockRepository) {
				f := make(url.Values)
				req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				expectingArgs := map[string]interface{}{}
				r.EXPECT().FindOffersByConditions(expectingArgs).Return(nil, errors.New("sample"))

				h := controllers.NewHandler(s, r)

				g.Expect(h.GetOffers(c)).ShouldNot(HaveOccurred())
				g.Expect(rec.Code).Should(Equal(http.StatusInternalServerError))

			},
		},
	}

	for _, c := range cases {
		s := mock_services.NewMockTaskService(mockCtrl)
		r := mock_repositories.NewMockRepository(mockCtrl)
		fmt.Println(c.description)
		c.expect(s, r)
		fmt.Println("ok")

	}
}
