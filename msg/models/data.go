package models

import (
	"time"
)

type SocketMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type StorageMessage struct {
	Id        int           `gorm:"primary_key, AUTO_INCREMENT"`
	TimeStamp string        `json:"timestamp"`
	Message   SocketMessage `json:"msg" gorm:"embedded"`
}

func CreateMessageSocket(user, msg string) *SocketMessage {
	return &SocketMessage{
		User:    user,
		Message: msg,
	}
}

func CreateMessageDB(msg *SocketMessage) *StorageMessage {
	if msg == nil {
		return nil
	}

	return &StorageMessage{
		TimeStamp: time.Now().Format(time.DateTime),
		Message:   *msg,
	}
}
