package db

import (
	"context"
	"encoding/json"
	"msg/models"

	"github.com/redis/go-redis/v9"
)

type WorkerDatabase struct {
	redis        *redis.Client
	redisListKey string
}

func NewWorkerDatabase(redisHost, redisPassword, redisListKey string) (*WorkerDatabase, error) {
	return &WorkerDatabase{
		redis: redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: redisPassword,
			DB:       0,
		}),
		redisListKey: redisListKey,
	}, nil
}

func (db *WorkerDatabase) Append(msg *models.StorageMessage) error {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = db.redis.ZAdd(context.Background(), db.redisListKey, redis.Z{
		Score:  float64(msg.TimeStamp.UnixNano()),
		Member: marshalled,
	}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (storage *WorkerDatabase) GetLatest(limit int) ([]models.StorageMessage, error) {
	res, err := storage.redis.ZRevRange(context.Background(), storage.redisListKey, 0, int64(limit)).Result()
	if err != nil {
		return nil, err
	}

	refSlice := make([]models.StorageMessage, min(limit, len(res)))
	for i := range refSlice {
		err := json.Unmarshal([]byte(res[i]), &refSlice[i])
		if err != nil {
			return nil, err
		}
	}

	return refSlice, nil
}
