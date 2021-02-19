package controllers

import (
	"MartellX/avito-tech-task/models"
	"MartellX/avito-tech-task/other"
	"MartellX/avito-tech-task/repository"
	"MartellX/avito-tech-task/services"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type Handler struct {
	Service *services.Service
	Repo    repository.Repository
}

func NewHandler(service *services.Service, repo repository.Repository) *Handler {
	return &Handler{Service: service, Repo: repo}
}

func (h *Handler) NewTask(ctx echo.Context) error {
	sellerId := ctx.FormValue("seller_id")
	url := ctx.FormValue("url")
	id, err := strconv.ParseUint(sellerId, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, other.GetJsonStatusMessage(http.StatusBadRequest, err.Error()))
	}

	task, err := h.Service.StartUploadingTask(uint(id), url)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, other.GetJsonStatusMessage(http.StatusBadRequest, err.Error()))
	}

	return ctx.JSON(task.StatusCode, task)
}

func (h *Handler) GetTask(ctx echo.Context) error {
	taskId := ctx.QueryParam("task_id")
	if taskId == "" {
		return ctx.JSON(http.StatusBadRequest,
			other.GetJsonStatusMessage(http.StatusBadRequest, "Не задан параметр task_id"))
	}

	task, ok := h.Service.GetTask(taskId)
	if !ok {
		return ctx.JSON(http.StatusNotFound,
			other.GetJsonStatusMessage(http.StatusNotFound, "Не найдено задание с таким id"))
	}
	return ctx.JSON(http.StatusOK, task)
}

func (h *Handler) GetOffers(ctx echo.Context) error {

	sellerIdStr := ctx.QueryParam("seller_id")
	offerIdStr := ctx.QueryParam("offer_id")
	name := ctx.QueryParam("name")

	// Используются аргументы:
	// seller_id uint
	// offer_id uint
	// name string
	args := map[string]interface{}{}

	if sellerIdStr != "" {
		sellerId, err := strconv.ParseUint(sellerIdStr, 10, 64)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				other.GetJsonStatusMessage(http.StatusBadRequest,
					fmt.Sprintf("Недопустимое значение для параметра %s, ожидалось %T", "seller_id", sellerId)))
		}
		args["seller_id"] = sellerId
	}

	if offerIdStr != "" {
		offerId, err := strconv.ParseUint(offerIdStr, 10, 64)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				other.GetJsonStatusMessage(http.StatusBadRequest, fmt.Sprintf("Недопустимое значение для параметра %s, ожидалось %T", "offer_id", offerId)))
		}
		args["offer_id"] = offerId
	}
	if name != "" {
		args["name"] = name
	}

	offers, err := h.Repo.FindOffersByConditions(args)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, other.GetJsonStatusMessage(http.StatusNotFound, err.Error()))
	}

	result := struct {
		Count int            `json:"count"`
		Items []models.Offer `json:"items"`
	}{len(offers), offers}

	return ctx.JSON(http.StatusOK, result)
}
