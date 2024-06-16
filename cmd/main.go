package main

import (
	"alert_bot/pkg/service"
	"path"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(path.Join("./", ".env")); err != nil {
		log.WithError(err).Warn("failed to load .env")
	}

	service := service.New()

	if err := service.Start(); err != nil {
		log.WithError(err).Fatal("service failed")
	}
}
