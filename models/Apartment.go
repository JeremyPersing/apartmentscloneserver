package models

import "gorm.io/gorm"

type Apartment struct {
	gorm.Model
	PropertyID uint    `json:"propertyID"`
	Unit       string  `json:"unit"`
	Bedrooms   int     `json:"bedrooms"`
	Bathrooms  float32 `json:"bathrooms"`
}
