package configs

import (
	"fmt"
	"msg/internal/db"
	"msg/models"

	"gorm.io/driver/postgres"
)

type DBConfig struct {
	Mock bool `json:"mock"`

	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     uint16 `json:"port"`
}

func (config DBConfig) GetStorage() (models.Storage, error) {
	if config.Mock {
		return &db.StorageMock{}, nil
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
		config.Host,
		config.User,
		config.Password,
		config.Port,
	)

	return db.NewGormDatabase(postgres.Open(dsn))
}

type ServerConfig struct {
	Address string `json:"address"`

	DB DBConfig `json:"db"`
}

func GetServerConfig(path string) (*ServerConfig, error) {
	return getConfig[ServerConfig](path)
}
