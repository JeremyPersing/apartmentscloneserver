package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Apartment struct {
	gorm.Model
	PropertyID  uint           `json:"propertyID"`
	Unit        string         `json:"unit"`
	Bedrooms    int            `json:"bedrooms"`
	Bathrooms   float32        `json:"bathrooms"`
	SqFt        int            `json:"sqFt"`
	Rent        float32        `json:"rent"`
	Deposit     float32        `json:"deposit"`
	LeaseLength string         `json:"leaseLength"`
	AvailableOn time.Time      `json:"availableOn"`
	Active      *bool          `json:"active"`
	Images      datatypes.JSON `json:"images"`
	Amenities   datatypes.JSON `json:"amenities"`
	Description string         `json:"description"`
}
