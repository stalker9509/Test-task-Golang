package service

import (
	"Test-task-Golang/internal/model/taskstruct"
	"Test-task-Golang/internal/service/taskorganize"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service interface {
	Create(task *taskstruct.Task) (string, error)
	Get(taskId string) (*taskstruct.TaskStatus, error)
}

type service struct {
	taskOrganize taskorganize.TaskOrganize
}

func NewService(taskOrganize taskorganize.TaskOrganize) *service {
	return &service{taskOrganize: taskOrganize}
}

func (service *service) Create(task *taskstruct.Task) (string, error) {
	return service.taskOrganize.Create(task)
}

func (service *service) Get(taskId string) (*taskstruct.TaskStatus, error) {
	return service.taskOrganize.Get(taskId)
}
