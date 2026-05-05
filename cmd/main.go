package main

import (
	"os"
	"strconv"

	dokkee "github.com/airvt1x/dokkee-backend"
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
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("failed to connect to postgres: %s", err.Error())
	}

	useSSL, _ := strconv.ParseBool(os.Getenv("S3_USE_SSL"))
	s3, err := repository.NewS3Repository(repository.S3Config{
		Endpoint:        os.Getenv("S3_ENDPOINT"),
		AccessKeyID:     os.Getenv("S3_ACCESS_KEY"),
		SecretAccessKey: os.Getenv("S3_SECRET_KEY"),
		BucketName:      os.Getenv("S3_BUCKET"),
		UseSSL:          useSSL,
	})
	if err != nil {
		logrus.Fatalf("failed to connect to S3: %s", err.Error())
	}

	repos := repository.NewRepository(db, s3)
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
