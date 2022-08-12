package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"

	"github.com/kataras/iris/v12"
)

func CreateProperty(ctx iris.Context) {
	var propertyInput PropertyInput

	err := ctx.ReadJSON(&propertyInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var apartments []models.Apartment
	bedroomLow := 0
	bedroomHigh := 0
	var bathroomLow float32 = 0.5
	var bathroomHigh float32 = 0.5

	for _, element := range propertyInput.Apartments {
		if element.Bathrooms < bathroomLow {
			bathroomLow = element.Bathrooms
		}
		if element.Bathrooms > bathroomHigh {
			bathroomHigh = element.Bathrooms
		}
		if *element.Bedrooms < bedroomLow {
			bedroomLow = *element.Bedrooms
		}
		if *element.Bedrooms > bedroomHigh {
			bedroomHigh = *element.Bedrooms
		}

		apartments = append(apartments, models.Apartment{
			Unit:      element.Unit,
			Bedrooms:  *element.Bedrooms,
			Bathrooms: element.Bathrooms,
		})
	}

	property := models.Property{
		UnitType:     propertyInput.UnitType,
		PropertyType: propertyInput.PropertyType,
		Street:       propertyInput.Street,
		City:         propertyInput.City,
		State:        propertyInput.State,
		Zip:          propertyInput.Zip,
		Lat:          propertyInput.Lat,
		Lng:          propertyInput.Lng,
		BedroomLow:   bedroomLow,
		BedroomHigh:  bedroomHigh,
		BathroomLow:  bathroomLow,
		BathroomHigh: bathroomHigh,
		Apartments:   apartments,
		ManagerID:    propertyInput.ManagerID,
	}

	storage.DB.Create(&property)

	ctx.JSON(property)
}

func GetProperty(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	var property models.Property
	propertyExists := storage.DB.Preload("Apartments").Find(&property, id)

	if propertyExists.Error != nil {
		utils.CreateError(
			iris.StatusInternalServerError,
			"Error", propertyExists.Error.Error(), ctx)
		return
	}

	if propertyExists.RowsAffected == 0 {
		utils.CreateError(iris.StatusNotFound, "Property Not Found", "Property Not Found", ctx)
		return
	}

	ctx.JSON(property)
}

type PropertyInput struct {
	UnitType     string           `json:"unitType" validate:"required,oneof=single multiple"`
	PropertyType string           `json:"propertyType" validate:"required,max=256"`
	Street       string           `json:"street" validate:"required,max=512"`
	City         string           `json:"city" validate:"required,max=512"`
	State        string           `json:"state" validate:"required,max=256"`
	Zip          int              `json:"zip" validate:"required"`
	Lat          float32          `json:"lat" validate:"required"`
	Lng          float32          `json:"lng" validate:"required"`
	ManagerID    uint             `json:"managerID" validate:"required"`
	Apartments   []ApartmentInput `json:"apartments" validate:"required,dive"`
}

type ApartmentInput struct {
	Unit      string  `json:"unit" validate:"max=512"`
	Bedrooms  *int    `json:"bedrooms" validate:"gte=0,max=6,required"` // make int a pointer so 0 is accepted
	Bathrooms float32 `json:"bathrooms" validate:"min=0.5,max=6.5,required"`
}
