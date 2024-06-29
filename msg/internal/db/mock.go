package db

import (
	"msg/models"
	"slices"
	"sort"
)

type StorageMock struct {
	messages []*models.StorageMessage
}

func (storage *StorageMock) Append(msg *models.StorageMessage) error {
	storage.messages = append(storage.messages, msg)

	return nil
}

func (storage *StorageMock) GetLatest(limit int) ([]models.StorageMessage, error) {
	sort.Slice(storage.messages, func(i, j int) bool {
		return storage.messages[i].TimeStamp.UnixMilli() > storage.messages[j].TimeStamp.UnixMilli()
	})

	refSlice := make([]models.StorageMessage, min(limit, len(storage.messages)))

	for i := 0; i < len(refSlice); i++ {
		refSlice[i] = *storage.messages[i]
	}
	slices.Reverse(refSlice)

	return refSlice, nil
}
