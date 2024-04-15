package main

import (
	"alert_bot/pkg/service"
	"path"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load(path.Join("./", ".env"))
	service := service.New()

	if err := service.Start(); err != nil {
		log.WithError(err).Fatal("service failed")
	}
}
