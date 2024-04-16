package storage

import "os"

type Storage interface {
	Subscribe(chatId int64) error
	Unsubscribe(chatId int64) error
	GetSubscribersUids() ([]int64, error)
}

func New() Storage {
	if _, ok := os.LookupEnv("REDIS_HOST"); !ok {
		return NewMemoryStorage()
	}

	return NewRedis()
}
