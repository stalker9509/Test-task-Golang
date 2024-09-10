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
	headerJSON, err := json.Marshal(task.Status.Headers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal headers: %w", err)
	}

	id := task.Status.ID

	query := "INSERT INTO task (ID ,Status, HTTPStatusCode, Headers, Length) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = repository.db.QueryRow(query, id, task.Status.Status, task.Status.HTTPStatusCode, headerJSON, task.Status.Length).
		Scan(&id)

	if err != nil {
		return "", err
	}

	return id, err
}

func (repository *Repository) Get(taskID string) (*model.TaskStatus, error) {
	var task model.TaskStatus
	var headerJSON []byte
	err := repository.db.QueryRow("SELECT ID, Status, HTTPStatusCode, Headers, Length FROM task WHERE id=$1", taskID).
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

func (repository *Repository) Update(task *model.Task) error {
	headerJSON, err := json.Marshal(task.Headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}
	query := `
		UPDATE task
		SET Status=$1, HTTPStatusCode=$2, Headers=$3, Length=$4;
`
	_, err = repository.db.Exec(query,
		task.Status.Status,
		task.Status.HTTPStatusCode,
		headerJSON,
		task.Status.Length)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}
