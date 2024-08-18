package service

import (
	"alert_bot/pkg/telegram"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	bot *telegram.Bot

	imageDir string

	logger *log.Entry
}

func New(logger *log.Entry) *Service {
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
		bot:      telegram.NewBot(logger),
		imageDir: imageDir,
		logger:   logger.WithField("type", "service"),
	}
}

func (s *Service) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.cleanData(); err != nil {
					return err
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	eg.Go(s.initRouter)
	eg.Go(s.bot.Start)

	return eg.Wait()
}

func (s *Service) initRouter() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /result", s.indexPost)
	mux.HandleFunc("POST /image", s.imagePost)
	mux.HandleFunc("GET /healthz", s.healthz)

	servicePort := fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT"))
	if servicePort == ":" {
		return fmt.Errorf("you should speciy port with SERVICE_PORT env")
	}

	s.logger.Infof("Service is listening on port %s", servicePort)

	return http.ListenAndServe(servicePort, mux)
}

func (s *Service) cleanData() error {
	dirEntry, err := os.ReadDir(s.imageDir)
	if err != nil {
		return fmt.Errorf("cannot read dir, err: %w", err)
	}

	for _, file := range dirEntry {
		if err := os.Remove(file.Name()); err != nil {
			s.logger.WithError(err).Errorf("cannot remove file %s", file.Name())
		}
	}

	return nil
}

func (s *Service) Close() {
	s.logger.Info("service was closed")
}
