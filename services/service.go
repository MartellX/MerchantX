package services

type TaskService interface {
	GetTask(id string) (*Task, bool)
	StartUploadingTask(sellerId uint64, xlsxURL string) (task *Task, err error)
}
