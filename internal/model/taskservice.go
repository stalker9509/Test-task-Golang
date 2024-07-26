package model

type TaskService interface {
	Create(task *Task) (string, error)
	Get(taskId string) (*TaskStatus, error)
}
