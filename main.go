package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Task struct {
	ID              string            `json:"id"`
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	Headers         map[string]string `json:"headers"`
	Status          string            `json:"status"`
	HTTPStatusCode  int               `json:"httpStatusCode"`
	ResponseHeaders http.Header       `json:"headers"`
	ResponseLength  int64             `json:"length"`
}

var taskStore = struct {
	sync.RWMutex
	tasks   map[string]Task
	counter int64
}{tasks: make(map[string]Task)}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskStore.Lock()
	defer taskStore.Unlock()
	task.ID = fmt.Sprintf("%d", taskStore.counter)
	taskStore.counter++
	task.Status = "new"
	taskStore.tasks[task.ID] = task

	go executeTask(taskStore.tasks[task.ID])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": task.ID})
}

func handleGetTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Path[len("/task/"):]

	taskStore.RLock()
	defer taskStore.RUnlock()
	task, ok := taskStore.tasks[taskID]
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func executeTask(task Task) {
	req, err := http.NewRequest(task.Method, task.URL, nil)
	if err != nil {
		task.Status = "error"
		return
	}

	for k, v := range task.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		task.Status = "error"
		return
	}
	defer resp.Body.Close()

	task.Status = "done"
	task.HTTPStatusCode = resp.StatusCode
	task.ResponseHeaders = resp.Header
	task.ResponseLength = resp.ContentLength

	taskStore.Lock()
	defer taskStore.Unlock()
	taskStore.tasks[task.ID] = task
}

func main() {
	http.HandleFunc("/task", handleCreateTask)
	http.HandleFunc("/task/", handleGetTaskStatus)

	fmt.Println("Starting server on :8080")
	fmt.Print(http.ListenAndServe(":8080", nil))
}
