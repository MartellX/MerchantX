package controllers

import (
	"MartellX/avito-tech-task/models"
	"MartellX/avito-tech-task/other"
	"MartellX/avito-tech-task/repositories"
	"MartellX/avito-tech-task/services"
	"fmt"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	TaskService services.TaskService
	Repo        repositories.Repository
}

func NewHandler(service services.TaskService, repo repositories.Repository) *Handler {
	return &Handler{TaskService: service, Repo: repo}
}

func (h *Handler) NewTask(ctx echo.Context) error {
	sellerId := ctx.FormValue("seller_id")
	url := ctx.FormValue("url")

	if sellerId == "" {
		return ctx.JSON(http.StatusBadRequest,
			other.GetJsonStatusMessage(http.StatusBadRequest, "Не задан параметр seller_id"))
	}
	if url == "" {
		return ctx.JSON(http.StatusBadRequest,
			other.GetJsonStatusMessage(http.StatusBadRequest, "Не задан параметр url"))
	}

	id, err := strconv.ParseUint(sellerId, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, other.GetJsonStatusMessage(http.StatusBadRequest, err.Error()))
	}

	task, err := h.TaskService.StartUploadingTask(id, url)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, other.GetJsonStatusMessage(http.StatusBadRequest, err.Error()))
	}

	return ctx.JSONPretty(task.StatusCode, task, "\t")
}

func (h *Handler) GetTask(ctx echo.Context) error {
	taskId := ctx.QueryParam("task_id")
	if taskId == "" {
		return ctx.JSON(http.StatusBadRequest,
			other.GetJsonStatusMessage(http.StatusBadRequest, "Не задан параметр task_id"))
	}

	task, ok := h.TaskService.GetTask(taskId)
	if !ok {
		return ctx.JSON(http.StatusNotFound,
			other.GetJsonStatusMessage(http.StatusNotFound, "Не найдено задание с таким id"))
	}
	return ctx.JSONPretty(http.StatusOK, task, "\t")
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
		if err == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, other.GetJsonStatusMessage(http.StatusNotFound, "Ничего не найдено"))
		} else {
			return ctx.JSON(http.StatusInternalServerError, "Непредвиденная ошибка")
		}

	}

	result := struct {
		Count int            `json:"count"`
		Items []models.Offer `json:"items"`
	}{len(offers), offers}

	return ctx.JSONPretty(http.StatusOK, result, "\t")
}
