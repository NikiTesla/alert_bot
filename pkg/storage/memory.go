package storage

type MemoryStorage struct {
	db map[int64]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		db: make(map[int64]string),
	}
}

func (mS *MemoryStorage) Subscribe(chatId int64) error {
	mS.db[chatId] = ""

	return nil
}

func (mS *MemoryStorage) Unsubscribe(chatId int64) error {
	delete(mS.db, chatId)

	return nil
}

func (mS *MemoryStorage) GetSubscribersUids() ([]int64, error) {
	var uids []int64

	for uid := range mS.db {
		uids = append(uids, uid)
	}

	return uids, nil
}
