package main

import (
	"os"

	"github.com/airvt1x/dokkee-backend"
	"github.com/airvt1x/dokkee-backend/internal/handler"
	"github.com/airvt1x/dokkee-backend/internal/repository"
	"github.com/airvt1x/dokkee-backend/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := godotenv.Load(); err != nil {
		logrus.Println("No .env file found")
	}

	if err := initConfig(); err != nil {
		logrus.Fatalf("failed to read config: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("failed to connect to postgres: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(dokkee.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("failed to start server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
