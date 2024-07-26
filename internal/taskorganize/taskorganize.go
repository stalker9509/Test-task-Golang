package taskorganize

import (
<<<<<<< HEAD:internal/taskorganize/taskorganize.go
	"Test-task-Golang/internal/taskstruct"
=======
	"Test-task-Golang/internal/model/taskstruct"
>>>>>>> parent of ae44e24 (upgrade):internal/service/taskorganize/taskorganize.go
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sync"
)

type TaskOrganize interface {
	Create(task *taskstruct.Task) (string, error)
	Get(taskId string) (*taskstruct.TaskStatus, error)
}

const (
	valueMaxGoroutine   = 5
	taskStatusInProcess = "in process"
	taskStatusError     = "error"
	taskStatusDone      = "done"
)

type organize struct {
	tasks        map[string]*taskstruct.Task
	taskQueue    chan *taskstruct.Task
	mutex        sync.Mutex
	WaitGroup    sync.WaitGroup
	maxGoroutine int
}

var ErrorTaskNotFound = errors.New("task not found")

func Init() *organize {
	return &organize{
		taskQueue:    make(chan *taskstruct.Task, valueMaxGoroutine),
		maxGoroutine: valueMaxGoroutine,
		tasks:        make(map[string]*taskstruct.Task),
	}
}

func (manager *organize) Create(task *taskstruct.Task) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	id := uuid.NewString()
	task.Status = &taskstruct.TaskStatus{
		ID:     id,
		Status: taskStatusInProcess,
	}
	manager.tasks[id] = task
	manager.taskQueue <- task
	manager.WaitGroup.Add(1)
	go manager.ExecuteTask()
	return id, nil
}

func (manager *organize) Get(taskID string) (*taskstruct.TaskStatus, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	task, ok := manager.tasks[taskID]
	if !ok {
		return nil, ErrorTaskNotFound
	}
	return task.Status, nil
}

func (manager *organize) ExecuteTask() {
	defer manager.WaitGroup.Done()
	for task := range manager.taskQueue {
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
		}
		task.Status.Length = len(body)
		task.Status.Status = taskStatusDone
	}
}
