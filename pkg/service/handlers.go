package service

import (
	"io"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

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

func (s *Service) healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK!\n"))
}
