package app

import (
	"Test-task-Golang/internal/handler"
	"Test-task-Golang/internal/service"
	"Test-task-Golang/internal/service/server"
	"Test-task-Golang/internal/service/taskorganize"
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	taskorganize := taskorganize.NewTaskOrganizeService()
	service := service.NewService(taskorganize)
	handler := handler.NewHandlerService(service)
	server := server.NewServer(handler.InitRout())

	go func() {
		err := server.Run()
		if err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()
	logrus.Print("Server Started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	taskorganize.WaitGroup.Wait()
	logrus.Print("Server Shutting Down")
	err := server.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}
