package service

import (
	"alert_bot/pkg/telegram"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	bot *telegram.Bot

	imageDir    string
	imagesSaved int
}

func New() *Service {
	imageDir := os.Getenv("IMAGES_DIR")
	if imageDir == "" {
		log.Warn("IMAGES_DIR env is empty, using ./")
		imageDir = "./"
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

func (s *Service) indexPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request's body", http.StatusBadRequest)
		return
	}

	if err := s.bot.NotifySubscribers(body); err != nil {
		http.Error(w, "Cannot notify subscribers", http.StatusInternalServerError)
		log.WithError(err).Error("cannot notify subscribers")
		return
	}
}

func (s *Service) imagePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Cannot read image from form. Did you send image in multipart form?", http.StatusBadRequest)
		log.WithError(err).Error("Cannot read image from form")
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read image data", http.StatusInternalServerError)
		log.WithError(err).Error("Could not read image data")
		return
	}

	imageName := path.Join(s.imageDir, fileHeader.Filename)
	if err := os.WriteFile(imageName, imageData, 0666); err != nil {
		http.Error(w, "Cannot save image", http.StatusInternalServerError)
		log.WithError(err).Error("Cannot save image")
		return
	}

	if err := s.bot.NotifySubscribersWithImage(imageName); err != nil {
		http.Error(w, "Cannot notify subscribers with image", http.StatusInternalServerError)
		log.WithError(err).Error("cannot notify subscribers with image")
		return
	}
}

func (s *Service) indexGet(w http.ResponseWriter, r *http.Request) {
	log.Info("get request was served")
}

func (s *Service) Close() {
	log.Info("service was closed")
}
