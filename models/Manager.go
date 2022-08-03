package models

import "gorm.io/gorm"

type Manager struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Website     string `json:"website"`
	UserID      uint   `json:"userID"`
	Image       string `json:"image"`
}
