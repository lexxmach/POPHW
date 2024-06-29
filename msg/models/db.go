package models

type Storage interface {
	GetLatest(limit int) ([]StorageMessage, error)
	Append(msg *StorageMessage) error
}
