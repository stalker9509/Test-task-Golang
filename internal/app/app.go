package app

import (
	"Test-task-Golang/internal/handler"
	"Test-task-Golang/internal/server"
	"Test-task-Golang/internal/service"
	"Test-task-Golang/internal/taskorganize"
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	taskorganize := taskorganize.Init()
	service := service.Init(taskorganize)
	handler := handler.Init(service)
	server := new(server.Server)
	go func() {
		err := server.Run(handler.InitRout())
		if err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
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
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
