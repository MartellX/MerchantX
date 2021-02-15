package services

import (
	"MartellX/avito-tech-task/models"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/gommon/log"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Service struct {
	repo  models.Repository
	tasks map[string]*Task
}

func NewService(repo models.Repository) *Service {
	return &Service{repo: repo, tasks: map[string]*Task{}}
}

type Task struct {
	TaskId     string `json:"task_id"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SellerId   uint   `json:"-"`

	Info struct {
		Created int `json:"created,omitempty"`
		Updated int `json:"updated,omitempty"`
		Deleted int `json:"deleted,omitempty"`
		Errors  int `json:"errors,omitempty"`
	} `json:"info,omitempty"`
}

func (t *Task) SetStatus(status string, code int) {
	t.Status = status
	t.StatusCode = code
}

func (s *Service) GetTask(id string) (*Task, bool) {
	task, ok := s.tasks[id]
	return task, ok
}

func (s *Service) createTask(sellerId uint) *Task {
	taskUUID, _ := uuid.DefaultGenerator.NewV4()
	id := taskUUID.String()
	task := &Task{
		TaskId:     id,
		Status:     "Created",
		StatusCode: http.StatusCreated,
		SellerId:   sellerId,
	}
	s.tasks[id] = task
	return task
}

func (s *Service) StartUploadingTask(sellerId uint, xlsxURL string) (task *Task, err error) {

	task = s.createTask(sellerId)

	go func() {
		req, err := http.Get(xlsxURL)
		if err != nil {
			log.Error(err)
			task.SetStatus(fmt.Sprintf("Error occured: %s", err), http.StatusBadRequest)
			return
		}

		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			task.SetStatus(fmt.Sprintf("Error occured: %s", err), http.StatusBadRequest)
			return
		}
		xlsxFile, err := xlsx.OpenBinary(body)
		if err != nil {
			task.SetStatus(fmt.Sprintf("Error occured: %s", err), http.StatusBadRequest)
			return
		}
		task.SetStatus("Parsing", http.StatusProcessing)
		ParsingTask(xlsxFile, task, s.repo)
	}()

	return task, nil
}

type RowData struct {
	Columns struct {
		OfferId   uint64 `xlsx:"0"`
		Name      string `xlsx:"1"`
		Price     int64  `xlsx:"2"`
		Quantity  int    `xlsx:"3"`
		Available bool   `xlsx:"4"`
	}
	ok  bool
	err error
}

func (r *RowData) UpdateColumns(offerId uint64, name string, price int64, quantity int, available bool) {
	r.Columns.OfferId = offerId
	r.Columns.Name = name
	r.Columns.Price = price
	r.Columns.Quantity = quantity
	r.Columns.Available = available
}

func ParsingTask(wb *xlsx.File, task *Task, repo models.Repository) {
	sh := wb.Sheets[0]

	rows := sh.Rows
	if len(rows) < 2 {
		task.SetStatus("Too few rows", http.StatusBadRequest)
		return
	}
	rows = sh.Rows[1:]
	parsedRows := make(chan RowData, len(rows))
	defer close(parsedRows)

	go checkAndUploadRows(parsedRows, task, repo)
	parsingRows(parsedRows, rows)
}

func parsingRows(parsedRows chan<- RowData, rows []*xlsx.Row) {
	for _, row := range rows {
		rowData := RowData{}
		rowData.ok = true
		cells := row.Cells
		if len(cells) >= 5 {
			offerIdStr, err := cells[0].GeneralNumericWithoutScientific()
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			offerId, err := strconv.ParseUint(offerIdStr, 10, 64)
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			name := cells[1].String()

			priceStr, err := cells[2].GeneralNumericWithoutScientific()
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			price, err := strconv.ParseInt(priceStr, 10, 64)
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			quantity, err := cells[3].Int()
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			availableStr, err := cells[4].FormattedValue()
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			available, err := strconv.ParseBool(availableStr)
			if err != nil {
				rowData.err = err
				rowData.ok = false
			}

			if rowData.ok {
				rowData.UpdateColumns(offerId, name, price, quantity, available)
			}

			if rowData.Columns.Price < 0 || rowData.Columns.Quantity < 0 {
				rowData.ok = false
			}

		} else {
			rowData.ok = false
		}

		parsedRows <- rowData
	}
}

func checkAndUploadRows(parsedRows <-chan RowData, task *Task, repo models.Repository) {
	defer task.SetStatus("Completed", http.StatusOK)
	sellerId := task.SellerId
	for parsedRow := range parsedRows {
		if !parsedRow.ok {
			task.Info.Errors++
			if parsedRow.err != nil {
				log.Debug(parsedRow.err)
			}
			continue
		}

		offerId := parsedRow.Columns.OfferId
		offer, err := repo.FindOffer(offerId, uint64(sellerId))
		if err == gorm.ErrRecordNotFound {
			if parsedRow.Columns.Available == false {
				continue
			}
			offer, err = repo.NewOffer(offerId, uint64(sellerId), parsedRow.Columns.Name, parsedRow.Columns.Price, parsedRow.Columns.Quantity, parsedRow.Columns.Available)
			if err != nil {
				task.Info.Errors++
				log.Debug(err)
				continue
			}
			task.Info.Created++
		} else if err == nil && offer != nil {
			if parsedRow.Columns.Available == false {
				repo.Delete(offer)
				task.Info.Deleted++
				continue
			}

			err := repo.UpdateColumns(offer, parsedRow.Columns.Name, parsedRow.Columns.Price, parsedRow.Columns.Quantity, parsedRow.Columns.Available)
			if err != nil {
				task.Info.Errors++
				log.Debug(err)
				continue
			}
			task.Info.Updated++
		} else {
			task.Info.Errors++
			if err != nil {
				log.Debug(err)
			}
		}
	}

}
