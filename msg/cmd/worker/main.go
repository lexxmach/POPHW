package main

import (
	"context"
	"encoding/json"
	"msg/internal"
	"msg/internal/configs"
	"msg/models"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	config := getConfig()
	storage := internal.Must(config.DB.GetStorage())

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.DB.Kafka.Host},
		Topic:    config.DB.Kafka.Topic,
		GroupID:  config.DB.Kafka.Topic,
		MaxBytes: 10e6, // 10MB
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		storageMsg := models.StorageMessage{}

		err = json.Unmarshal(msg.Value, &storageMsg)
		logger.Debugf(
			"Got message from user %q: %q",
			storageMsg.Message.User,
			storageMsg.Message.Message,
		)
		if err != nil {
			panic(err)
		}

		err = storage.Append(&storageMsg)
		if err != nil {
			panic(err)
		}

		err = r.CommitMessages(context.Background(), msg)
		if err != nil {
			panic(err)
		}
	}
}

func getConfig() *configs.WorkerConfig {
	configPath := pflag.StringP("config", "c", "", "Path to config")
	pflag.Parse()

	return internal.Must(configs.GetWorkerConfig(*configPath))
}
