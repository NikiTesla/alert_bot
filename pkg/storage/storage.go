package storage

import (
	"alert_bot/pkg/model"
	"fmt"
	"os"
)

type Storage interface {
	Subscribe(chatId int64) error
	Unsubscribe(chatId int64) error
	GetSubscribersUids() ([]int64, error)
	SetStatus(chatId int64, status model.Status) error
	GetStatus(chatId int64) (model.Status, error)
}

func New() (Storage, error) {
	if os.Getenv("IN_MEMORY") == "true" {
		return NewMemoryStorage(), nil
	}

	storageFile := os.Getenv("DB_FILENAME")
	if storageFile == "" {
		storageFile = sqliteStorageFilename
	}
	storage, err := NewSQLiteStorage(storageFile)
	if err != nil {
		return nil, fmt.Errorf("creating sqlte storage, err: %w", err)
	}
	return storage, nil
}
