package services

import (
	"MartellX/avito-tech-task/models"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/gommon/log"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
)

var tasks = map[string] *Task{}
var tasksLock = sync.Mutex{}

type Task struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
	StatusCode int `json:"status_code"`
	SellerId uint `json:"-"`

	Info struct{
		Created int	`json:"created,omitempty"`
		Updated int	`json:"updated,omitempty"`
		Deleted int `json:"deleted,omitempty"`
		Errors  int `json:"errors,omitempty"`
	} `json:"info,omitempty"`
}

func (t *Task) SetStatus(status string, code int) {
	t.Status = status
	t.StatusCode = code
}

func GetTask(id string) (*Task, bool) {
	task, ok := tasks[id]
	return task, ok
}

func CreateTask(sellerId uint) *Task {
	taskUUID, _ := uuid.DefaultGenerator.NewV4()
	id := taskUUID.String()
	task := &Task{
		TaskId:     id,
		Status:     "Created",
		StatusCode: http.StatusCreated,
		SellerId: sellerId,
	}
	tasks[id] = task
	return task
}

func StartUploadingTask(sellerId uint, xlsxURL string) (task *Task, err error) {
	req, err := http.Get(xlsxURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	task = CreateTask(sellerId)

	go func() {
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
		ParsingTask(xlsxFile, task)
	}()

	return task, nil
}

type RowData struct {
	Columns struct{
		OfferId int	`xlsx:"0"`
		Name string		`xlsx:"1"`
		Price float64	`xlsx:"2"`
		Quantity int	`xlsx:"3"`
		Available bool	`xlsx:"4"`
	}
	ok bool
	err error
}

func (r *RowData) UpdateColumns(offerId int, name string, price float64, quantity int, available bool) {
	r.Columns.OfferId = offerId
	r.Columns.Name = name
	r.Columns.Price = price
	r.Columns.Quantity = quantity
	r.Columns.Available = available
}

func ParsingTask(wb *xlsx.File, task *Task) {
	sh := wb.Sheets[0]

	rows := sh.Rows
	if len(rows) < 2 {
		task.SetStatus("Too few rows", http.StatusBadRequest)
		return
	}
	rows = sh.Rows[1:]
	parsedRows := make(chan RowData, len(rows))
	defer close(parsedRows)

	go checkAndUploadRows(parsedRows, task)
	parsingRows(parsedRows, rows)
}

func parsingRows(parsedRows chan<- RowData, rows []*xlsx.Row)  {
	for _, row := range rows {
		rowData := RowData{}
		cells := row.Cells
		if len(cells) == 5 {
			offerId, err := cells[0].Int()
			name := cells[1].String()
			price, err := cells[2].Float()
			quantity, err := cells[3].Int()
			available := cells[4].Bool()

			if err != nil {
				rowData.err = err
				rowData.ok = false
			} else {
				rowData.UpdateColumns(offerId, name, price, quantity, available)
				rowData.ok = true
				if rowData.Columns.OfferId < 0 || rowData.Columns.Price < 0 || rowData.Columns.Quantity < 0 || rowData.Columns.Price == math.NaN() {
					rowData.ok = false
				}
			}
		}
		//if err != nil {
		//	rowData.err = err
		//	rowData.ok = false
		//} else {
		//	rowData.ok = true
		//}
		//if rowData.Columns.OfferId < 0 || rowData.Columns.Price < 0 || rowData.Columns.Quantity < 0 {
		//	rowData.ok = false
		//}
		parsedRows <- rowData
	}

}

func checkAndUploadRows(parsedRows <-chan RowData, task *Task) {
	sellerId := task.SellerId
	updated := map[uint] int {}
	for parsedRow := range parsedRows {
		if !parsedRow.ok {
			task.Info.Errors++
			log.Info(parsedRow.err)
			continue
		}

		offerId := parsedRow.Columns.OfferId
		offer, err := models.FindOffer(uint(offerId), sellerId)
		if err == gorm.ErrRecordNotFound {
			if parsedRow.Columns.Available == false {
				continue
			}
			offer, err = models.NewOffer(uint(offerId), sellerId, parsedRow.Columns.Name, parsedRow.Columns.Price, parsedRow.Columns.Quantity, parsedRow.Columns.Available)
			if err != nil {
				task.Info.Errors++
				log.Info(err)
				continue
			}
			task.Info.Created++
		} else if err == nil && offer != nil{
			if parsedRow.Columns.Available == false {
				offer.Delete()
				task.Info.Deleted++
				continue
			}

			err := offer.UpdateColumns(parsedRow.Columns.Name, parsedRow.Columns.Price, parsedRow.Columns.Quantity, parsedRow.Columns.Available)
			if err != nil {
				task.Info.Errors++
				log.Info(err)
				continue
			}
			task.Info.Updated++
			updated[offer.OfferId]++
		} else {
			task.Info.Errors++
			if err != nil {
				log.Info(err)
			}
		}
	}
	task.SetStatus("Completed", http.StatusOK)

	fmt.Printf("Overupdates: %v \n", updated)
}