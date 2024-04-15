package main

import (
	"alert_bot/pkg/service"

	log "github.com/sirupsen/logrus"
)

func main() {
	service := service.New()

	if err := service.Start(); err != nil {
		log.WithError(err).Fatal("service failed")
	}
}
