package storage

import (
	"alert_bot/pkg/errs"
	"alert_bot/pkg/model"
	"sync"
)

type MemoryStorage struct {
	chatIdToStatusDb map[int64]model.Status
	mx               sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		chatIdToStatusDb: make(map[int64]model.Status),
		mx:               sync.RWMutex{},
	}
}

func (mS *MemoryStorage) Subscribe(chatId int64) error {
	mS.mx.Lock()
	defer mS.mx.Unlock()
	mS.chatIdToStatusDb[chatId] = ""

	return nil
}

func (mS *MemoryStorage) Unsubscribe(chatId int64) error {
	mS.mx.Lock()
	defer mS.mx.Unlock()
	delete(mS.chatIdToStatusDb, chatId)

	return nil
}

func (mS *MemoryStorage) GetSubscribersUids() ([]int64, error) {
	mS.mx.RLock()
	defer mS.mx.RUnlock()

	uids := make([]int64, 0, len(mS.chatIdToStatusDb))
	for uid := range mS.chatIdToStatusDb {
		uids = append(uids, uid)
	}

	return uids, nil
}

func (mS *MemoryStorage) SetStatus(chatId int64, status model.Status) error {
	mS.mx.Lock()
	defer mS.mx.Unlock()
	mS.chatIdToStatusDb[chatId] = status

	return nil
}

func (mS *MemoryStorage) GetStatus(chatId int64) (model.Status, error) {
	mS.mx.RLock()
	defer mS.mx.RUnlock()
	status, ok := mS.chatIdToStatusDb[chatId]
	if !ok {
		return "", errs.ErrNotFound
	}
	return status, nil
}
