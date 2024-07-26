package service

import (
<<<<<<< HEAD
	"Test-task-Golang/internal/taskorganize"
	"Test-task-Golang/internal/taskstruct"
=======
	"Test-task-Golang/internal/model/taskstruct"
	"Test-task-Golang/internal/service/taskorganize"
>>>>>>> parent of ae44e24 (upgrade)
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service interface {
	Create(task *taskstruct.Task) (string, error)
	Get(taskId string) (*taskstruct.TaskStatus, error)
}

type service struct {
	taskOrganize taskorganize.TaskOrganize
}

<<<<<<< HEAD
func Init(taskOrganize taskorganize.TaskOrganize) *service {
=======
func NewService(taskOrganize taskorganize.TaskOrganize) *service {
>>>>>>> parent of ae44e24 (upgrade)
	return &service{taskOrganize: taskOrganize}
}

func (service *service) Create(task *taskstruct.Task) (string, error) {
	return service.taskOrganize.Create(task)
}

func (service *service) Get(taskId string) (*taskstruct.TaskStatus, error) {
	return service.taskOrganize.Get(taskId)
}
