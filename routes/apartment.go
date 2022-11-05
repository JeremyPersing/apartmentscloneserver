package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func GetApartmentsByPropertyID(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	var apartments []models.Apartment
	apartmentsExist := storage.DB.Where("property_id = ?", id).Find(&apartments)

	if apartmentsExist.Error != nil {
		utils.CreateError(
			iris.StatusInternalServerError,
			"Error", apartmentsExist.Error.Error(), ctx)
		return
	}

	ctx.JSON(apartments)
}

func UpdateApartments(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	property := GetPropertyAndAssociationsByPropertyID(id, ctx)
	if property == nil {
		return
	}

	claims := jwt.Get(ctx).(*utils.AccessToken)

	if property.UserID != claims.ID {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	var updatedApartments []UpdateUnitsInput
	err := ctx.ReadJSON(&updatedApartments)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var newApartments []models.Apartment
	bedroomLow := property.BedroomLow
	bedroomHigh := property.BedroomHigh
	var bathroomLow float32 = property.BathroomLow
	var bathroomHigh float32 = property.BathroomHigh

	for _, apartment := range updatedApartments {
		if *apartment.Bedrooms > bedroomHigh {
			bedroomHigh = *apartment.Bedrooms
		}
		if *apartment.Bedrooms < bedroomLow {
			bedroomLow = *apartment.Bedrooms
		}
		if apartment.Bathrooms > bathroomHigh {
			bathroomHigh = apartment.Bathrooms
		}
		if apartment.Bathrooms < bathroomLow {
			bathroomLow = apartment.Bathrooms
		}

		currApartment := models.Apartment{
			Unit:        apartment.Unit,
			Bedrooms:    *apartment.Bedrooms,
			Bathrooms:   apartment.Bathrooms,
			SqFt:        apartment.SqFt,
			Active:      apartment.Active,
			AvailableOn: apartment.AvailableOn,
			PropertyID:  property.ID,
		}

		if apartment.ID != nil {
			currApartment.ID = *apartment.ID
			storage.DB.Model(&currApartment).Updates(currApartment)
		} else {
			newApartments = append(newApartments, currApartment)
		}
	}

	if len(newApartments) > 0 {
		rowsUpdated := storage.DB.Create(&newApartments)

		if rowsUpdated.Error != nil {
			utils.CreateError(
				iris.StatusInternalServerError,
				"Error", rowsUpdated.Error.Error(), ctx)
			return
		}
	}

	ctx.StatusCode(iris.StatusNoContent)
}

type UpdateUnitsInput struct {
	ID          *uint     `json:"ID"`
	Unit        string    `json:"unit" validate:"max=512"`
	Bedrooms    *int      `json:"bedrooms" validate:"gte=0,max=6,required"` // make int a pointer so 0 is accepted
	Bathrooms   float32   `json:"bathrooms" validate:"min=0.5,max=6.5,required"`
	SqFt        int       `json:"sqFt" validate:"max=100000000000,required"`
	Active      *bool     `json:"active" validate:"required"`
	AvailableOn time.Time `json:"availableOn" validate:"required"`
}
