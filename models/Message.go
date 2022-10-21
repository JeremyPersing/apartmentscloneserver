package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ConversationID uint
	SenderID       uint   `json:"senderID"`
	ReceiverID     uint   `json:"receiverID"`
	Text           string `json:"text"`
}
