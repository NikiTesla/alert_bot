package main

import (
	"alert_bot/pkg/service"
	"context"
	"os"
	"path"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(path.Join("./", ".env")); err != nil {
		log.WithError(err).Warn("failed to load .env")
	}

	service := service.New(createLogger())

	if err := service.Start(context.Background()); err != nil {
		log.WithError(err).Fatal("service failed")
	}
}

func createLogger() *log.Entry {
	logger := log.NewEntry(log.StandardLogger())
	if os.Getenv("DEBUG") == "true" {
		logger.Logger.SetLevel(log.DebugLevel)
	}
	return logger
}
