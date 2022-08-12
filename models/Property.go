package models

import "gorm.io/gorm"

type Property struct {
	gorm.Model
	UnitType     string      `json:"unitType"`
	PropertyType string      `json:"propertyType"`
	Street       string      `json:"street"`
	City         string      `json:"city"`
	State        string      `json:"state"`
	Zip          int         `json:"zip"`
	Lat          float32     `json:"lat"`
	Lng          float32     `json:"lng"`
	BedroomLow   int         `json:"bedroomLow"`   // calculate based off apartments
	BedroomHigh  int         `json:"bedroomHigh"`  // calculate based off apartments
	BathroomLow  float32     `json:"bathroomLow"`  // calculate based off apartments
	BathroomHigh float32     `json:"bathroomHigh"` // calculate based off apartments
	ManagerID    uint        `json:"managerID"`
	Apartments   []Apartment `json:"apartments"`
}
