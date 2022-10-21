package models

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	TenantID   uint      `json:"tenantID"`
	OwnerID    uint      `json:"ownerID"`
	PropertyID uint      `json:"propertyID"`
	Messages   []Message `json:"messages"`
}
