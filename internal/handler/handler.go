package handler

import (
	"Test-task-Golang/internal/model"
	"Test-task-Golang/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service model.TaskService
}

func NewHandlerService() *Handler {
	return InitService(service.NewService())
}

func InitService(service model.TaskService) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) InitRout() *gin.Engine {
	rout := gin.Default()
	rout.POST("/task", handler.CreatTask)
	rout.GET("/task/:id", handler.GetTask)
	return rout
}

func (handler *Handler) CreatTask(ctx *gin.Context) {
	var task model.Task
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := handler.service.Create(&task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (handler *Handler) GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := handler.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}
