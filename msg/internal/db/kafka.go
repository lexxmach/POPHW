package db

import (
	"context"
	"encoding/json"
	"msg/models"
	"slices"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type KafkaDatabase struct {
	writer *kafka.Writer
	redis  *redis.Client

	redisListKey string
}

func NewKafkaDatabase(redisHost, redisPassword, redisListKey, kafkaHost, topic string) (models.Storage, error) {
	return &KafkaDatabase{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(kafkaHost),
			Topic: topic,
		},
		redis: redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: redisPassword,
			DB:       0,
		}),
		redisListKey: redisListKey,
	}, nil
}

func (storage *KafkaDatabase) Append(msg *models.StorageMessage) error {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = storage.writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(marshalled),
	})
	if err != nil {
		return err
	}

	return nil
}

func (storage *KafkaDatabase) GetLatest(limit int) ([]models.StorageMessage, error) {
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
	slices.Reverse(refSlice)

	return refSlice, nil
}
