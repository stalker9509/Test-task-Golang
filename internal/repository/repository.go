package repository

import (
	"Test-task-Golang/internal/model"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) Create(task *model.Task) (string, error) {
	headerJSON, err := json.Marshal(task.Headers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal headers: %w", err)
	}

	id := task.Status.ID

	query := "INSERT INTO tasks (Status, HTTPStatusCode, Headers, Length) VALUES ($1, $2, $3, $4) RETURNING id"
	err = repository.db.QueryRow(query, task.Status.Status, task.Status.HTTPStatusCode, headerJSON, task.Status.Length).
		Scan(&id)

	if err != nil {
		return "", err
	}

	return id, err
}

func (repository *Repository) Get(taskID string) (*model.TaskStatus, error) {
	var task model.TaskStatus
	var headerJSON []byte
	err := repository.db.QueryRow("SELECT ID, Status, HTTPStatusCode, Headers, Length FROM tasks WHERE id=$1", taskID).
		Scan(&task.ID, &task.Status, &task.HTTPStatusCode, &headerJSON, &task.Length)
	if err != nil {
		return &model.TaskStatus{}, err
	}

	err = json.Unmarshal(headerJSON, &task.Headers)
	if err != nil {
		return &model.TaskStatus{}, err
	}

	return &task, nil
}
