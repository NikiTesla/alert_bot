package storage

import "alert_bot/pkg/model"

type Storage interface {
	Subscribe(chatId int64) error
	Unsubscribe(chatId int64) error
	GetSubscribersUids() ([]int64, error)
	SetStatus(chatId int64, status model.Status) error
	GetStatus(chatId int64) (model.Status, error)
}

func New() Storage {
	return NewMemoryStorage()
}
