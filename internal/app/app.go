package app

import (
	"Test-task-Golang/internal/handler"
	"Test-task-Golang/internal/service"
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	service := service.NewService()
	handler := handler.NewHandlerService(service)
	server := NewServer(handler.InitRout())

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
	logrus.Print("Server Shutting Down")
	err := server.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}
