package app

import (
	"Test-task-Golang/internal/handler"
	"Test-task-Golang/internal/repository"
	"Test-task-Golang/internal/service"
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	err := InitConfig()
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Configs{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
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

	logrus.Info("Task service shutting down")
	err = service.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error occurred on service shutting down: %s", err.Error())
	}

	logrus.Print("Server Shutting Down")
	err = server.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
