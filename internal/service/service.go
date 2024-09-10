package service

import (
	"Test-task-Golang/internal/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

const (
	valueMaxGoroutine   = 5
	taskStatusInProcess = "in process"
	taskStatusError     = "error"
	taskStatusDone      = "done"
	httpClientTimeout   = 10 * time.Second
)

type Service struct {
	tasks        map[string]*model.Task
	taskQueue    chan *model.Task
	mutex        sync.RWMutex
	httpClient   *http.Client
	maxGoroutine int
	repository   model.TaskService
}

var ErrorTaskNotFound = errors.New("task not found")

func NewService(repository model.TaskService) *Service {
	manager := &Service{
		taskQueue:    make(chan *model.Task, valueMaxGoroutine),
		maxGoroutine: valueMaxGoroutine,
		tasks:        make(map[string]*model.Task),
		httpClient:   &http.Client{Timeout: httpClientTimeout},
		repository:   repository,
	}
	manager.startWorkers()

	return manager
}

func (manager *Service) startWorkers() {
	for i := 0; i < manager.maxGoroutine; i++ {
		go manager.worker()
	}
}

func (manager *Service) worker() {
	for task := range manager.taskQueue {
		manager.executeTask(task)
	}
}

func (manager *Service) Shutdown(ctx context.Context) error {
	close(manager.taskQueue)
	select {
	case <-manager.taskQueue:
		logrus.Info("All workers finished")
		return nil
	case <-ctx.Done():
		logrus.Warn("Shutdown was not successful, not all workers finished")
		return ctx.Err()
	}
}

func (manager *Service) Create(task *model.Task) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	id := uuid.NewString()
	task.Status = &model.TaskStatus{
		ID:     id,
		Status: taskStatusInProcess,
	}
	manager.tasks[id] = task
	id, err := manager.repository.Create(task)
	if err != nil {
		return "", err
	}
	manager.taskQueue <- task
	return id, nil
}

func (manager *Service) Get(taskID string) (*model.TaskStatus, error) {
	return manager.repository.Get(taskID)
}

func (manager *Service) executeTask(task *model.Task) {
	request, err := http.NewRequest(task.Method, task.URL, nil)
	if err != nil {
		task.Status.Status = taskStatusError
		manager.updateTaskStatus(task)
		return
	}
	for key, value := range task.Headers {
		request.Header.Set(key, value)
	}
	response, err := manager.httpClient.Do(request)
	if err != nil {
		task.Status.Status = taskStatusError
		manager.updateTaskStatus(task)
		return
	}
	defer response.Body.Close()
	task.Status.HTTPStatusCode = response.StatusCode
	headers := make(map[string]string)
	for key, value := range response.Header {
		headers[key] = value[0]
	}
	task.Status.Headers = headers
	body, err := io.ReadAll(response.Body)
	if err != nil {
		task.Status.Status = taskStatusError
		manager.updateTaskStatus(task)
		return
	}
	task.Status.Length = len(body)
	task.Status.Status = taskStatusDone
	manager.updateTaskStatus(task)
}

func (manager *Service) Update(task *model.Task) error {
	err := manager.repository.Update(task)
	if err != nil {
		logrus.Errorf("Failed to update task: %v", err)
		return err
	}
	return nil
}

func (manager *Service) updateTaskStatus(task *model.Task) {
	err := manager.repository.Update(task)
	if err != nil {
		logrus.Errorf("Failed to update task status: %v", err)
	}
}
