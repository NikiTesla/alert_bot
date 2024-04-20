package service

import (
	"alert_bot/pkg/telegram"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	bot *telegram.Bot

	imageDir string
}

func New() *Service {
	imageDir := os.Getenv("IMAGES_DIR")
	if imageDir == "" {
		log.Warn("IMAGES_DIR env is empty, using ./images")
		imageDir = "./images"
	}

	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.WithError(err).Fatalf("could not create dir %s", imageDir)
		}
	}

	return &Service{
		bot:      telegram.NewBot(),
		imageDir: imageDir,
	}
}

func (s *Service) Start() error {
	go func() {
		if err := s.initRouter(); err != nil {
			log.WithError(err).Fatal("router has failed")
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Hour)
		for range ticker.C {
			s.cleanData()
		}
	}()

	return s.bot.Start()
}

func (s *Service) initRouter() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /result", s.indexPost)
	mux.HandleFunc("POST /image", s.imagePost)
	mux.HandleFunc("GET /result", s.indexGet)

	servicePort := fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT"))
	log.Infof("Service is listening on port %s", servicePort)

	return http.ListenAndServe(servicePort, mux)
}

func (s *Service) cleanData() error {
	dirEntry, err := os.ReadDir(s.imageDir)
	if err != nil {
		return fmt.Errorf("cannot read dir, err: %w", err)
	}

	for _, file := range dirEntry {
		if err := os.Remove(file.Name()); err != nil {
			log.WithError(err).Errorf("cannot remove file %s", file.Name())
		}
	}

	return nil
}

func (s *Service) Close() {
	log.Info("service was closed")
}
