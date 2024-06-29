package configs

import (
	"encoding/json"
	"fmt"
	"msg/internal/db"
	"msg/models"
	"os"

	"gorm.io/driver/postgres"
)

type KafkaConfig struct {
	Host  string `json:"host"`
	Topic string `json:"topic"`
}

type DBConfig struct {
	// mock/postgres/kafka/worker
	Type string `json:"type"`

	Host     string `json:"host,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Port     uint16 `json:"port,omitempty"`

	Kafka *KafkaConfig `json:"kafka,omitempty"`

	// Redis
	ListKey string `json:"list_key,omitempty"`
}

func (config DBConfig) GetStorage() (models.Storage, error) {
	if config.Type == "mock" {
		return &db.StorageMock{}, nil
	} else if config.Type == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
			config.Host,
			config.User,
			config.Password,
			config.Port,
		)

		return db.NewGormDatabase(postgres.Open(dsn))
	} else if config.Type == "kafka" {
		return db.NewKafkaDatabase(
			config.Host,
			config.Password,
			config.ListKey,
			config.Kafka.Host,
			config.Kafka.Topic,
		)
	} else if config.Type == "worker" {
		return db.NewWorkerDatabase(
			config.Host,
			config.Password,
			config.ListKey,
		)
	} else {
		return nil, fmt.Errorf("unsupported db type %q", config.Type)
	}
}

func getConfig[T any](path string) (*T, error) {
	plan, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(T)
	err = json.Unmarshal(plan, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
