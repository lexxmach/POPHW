package db

import (
	"fmt"
	"msg/models"

	"gorm.io/gorm"
)

type GormDatabase struct {
	db *gorm.DB
}

func NewGormDatabase(dialector gorm.Dialector) (models.Storage, error) {
	gormDB, err := gorm.Open(dialector)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gorm: %w", err)
	}

	gormDB.AutoMigrate(&models.StorageMessage{})

	return &GormDatabase{db: gormDB}, nil
}

func (g *GormDatabase) Append(msg *models.StorageMessage) error {
	tx := g.db.Create(msg)
	return tx.Error
}

func (g *GormDatabase) GetLatest(limit int) ([]models.StorageMessage, error) {
	messages := make([]models.StorageMessage, 0)
	tx := g.db.Order("time_stamp desc").Limit(limit).Find(&messages)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}
