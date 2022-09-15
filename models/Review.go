package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID     uint   `json:"userID"`
	PropertyID uint   `json:"propertyID"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Stars      int    `json:"stars"`
}
