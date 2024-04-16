package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	redisOpts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	client := redis.NewClient(redisOpts)
	if cmd := client.Ping(context.Background()); cmd != nil && cmd.Err() != nil {
		log.WithError(cmd.Err()).Fatal("failed to ping database")
	}
	log.Info("redis client was created")

	return &Redis{
		client: client,
	}
}

func (r *Redis) Subscribe(chatId int64) error {
	// r.db[chatId] = ""

	return nil
}

func (r *Redis) Unsubscribe(chatId int64) error {
	// delete(r.db, chatId)

	return nil
}

func (r *Redis) GetSubscribersUids() ([]int64, error) {
	var uids []int64

	// for uid := range r.db {
	// 	uids = append(uids, uid)
	// }

	return uids, nil
}
