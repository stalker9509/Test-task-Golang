package service

import (
	"Test-task-Golang/internal/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type service struct {
	taskOrganize model.TaskService
}

func NewService(taskOrganize model.TaskService) *service {
	return &service{taskOrganize: taskOrganize}
}

func (service *service) Create(task *model.Task) (string, error) {
	return service.taskOrganize.Create(task)
}

func (service *service) Get(taskId string) (*model.TaskStatus, error) {
	return service.taskOrganize.Get(taskId)
}
