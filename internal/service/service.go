package service

import (
	"Test-task-Golang/internal/model"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

const (
	valueMaxGoroutine   = 5
	taskStatusInProcess = "in process"
	taskStatusError     = "error"
	taskStatusDone      = "done"
)

type Service struct {
	tasks        map[string]*model.Task
	taskQueue    chan *model.Task
	mutex        sync.Mutex
	WaitGroup    sync.WaitGroup
	maxGoroutine int
}

var ErrorTaskNotFound = errors.New("task not found")

func NewService() *Service {
	manager := &Service{
		taskQueue:    make(chan *model.Task, valueMaxGoroutine),
		maxGoroutine: valueMaxGoroutine,
		tasks:        make(map[string]*model.Task),
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
		manager.ExecuteTask(task)
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
	manager.taskQueue <- task
	return id, nil
}

func (manager *Service) Get(taskID string) (*model.TaskStatus, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	task, ok := manager.tasks[taskID]
	if !ok {
		return nil, ErrorTaskNotFound
	}
	return task.Status, nil
}

func (manager *Service) ExecuteTask(task *model.Task) {
	client := &http.Client{}
	request, err := http.NewRequest(task.Method, task.URL, nil)
	if err != nil {
		task.Status.Status = taskStatusError
		return
	}
	for key, value := range task.Headers {
		request.Header.Set(key, value)
	}
	response, err := client.Do(request)
	if err != nil {
		task.Status.Status = taskStatusError
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
		return
	}
	task.Status.Length = len(body)
	task.Status.Status = taskStatusDone
}
