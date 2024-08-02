package handler

import (
	"Test-task-Golang/internal/handler"
	"Test-task-Golang/internal/model"
	mockservice "Test-task-Golang/internal/service/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockservice.NewMockService(ctrl)
	handler := handler.NewHandlerService(mockService)

	router := gin.Default()
	router.POST("/task", handler.CreatTask)

	t.Run("Create task successfully", func(t *testing.T) {
		task := &model.Task{Method: "GET", URL: "http://example.com"}
		mockService.EXPECT().Create(task).Return("1", nil)

		body := strings.NewReader(`{"method":"GET","url":"http://example.com"}`)
		req, _ := http.NewRequest("POST", "/task", body)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		expected := `{"id":"1"}`

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expected, resp.Body.String())
	})

	t.Run("Create task with invalid JSON body", func(t *testing.T) {
		body := strings.NewReader(`{"method":"GET","url":"http://example.com"`)
		req, _ := http.NewRequest("POST", "/task", body)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		expected := `{"error":"unexpected EOF"}`

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, expected, resp.Body.String())
	})

	t.Run("Create task with service error", func(t *testing.T) {
		task := &model.Task{Method: "GET", URL: "http://example.com"}
		mockService.EXPECT().Create(task).Return("", errors.New("service error"))

		body := strings.NewReader(`{"method":"GET","url":"http://example.com"}`)
		req, _ := http.NewRequest("POST", "/task", body)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		expected := `{"error":"service error"}`

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, expected, resp.Body.String())
	})
}

func TestGetTaskStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockservice.NewMockService(ctrl)
	handler := handler.NewHandlerService(mockService)

	router := gin.Default()
	router.GET("/task/:id", handler.GetTask)

	t.Run("Get task status successfully", func(t *testing.T) {
		taskStatus := &model.TaskStatus{ID: "1", Status: "done", HTTPStatusCode: 200}
		mockService.EXPECT().Get("1").Return(taskStatus, nil)

		req, _ := http.NewRequest("GET", "/task/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		expected := `{"id":"1","status":"done","httpStatusCode":200,"headers":null,"length":0}`

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expected, resp.Body.String())
	})

	t.Run("Get task status with service error", func(t *testing.T) {
		mockService.EXPECT().Get("1").Return(nil, errors.New("service error"))

		req, _ := http.NewRequest("GET", "/task/1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		expected := `{"error":"service error"}`

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, expected, resp.Body.String())
	})
}
