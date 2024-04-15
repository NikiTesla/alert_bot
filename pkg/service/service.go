package service

import (
	"alert_bot/pkg/telegram"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	bot *telegram.Bot
}

func New() *Service {
	return &Service{
		bot: telegram.NewBot(),
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

	mux.HandleFunc("POST /", s.indexPost)
	mux.HandleFunc("GET /", s.indexGet)

	servicePort := ":2704"
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

	s.bot.NotifySubscribers(body)
}

func (s *Service) indexGet(w http.ResponseWriter, r *http.Request) {
	log.Info("get request was served")
}

func (s *Service) Close() {
	log.Info("service was closed")
}
