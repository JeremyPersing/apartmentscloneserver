package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm/clause"
)

func CreateProperty(ctx iris.Context) {
	var propertyInput CreatePropertyInput

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
			Unit:        element.Unit,
			Bedrooms:    *element.Bedrooms,
			Bathrooms:   element.Bathrooms,
			AvailableOn: element.AvailableOn,
			Active:      element.Active,
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
		UserID:       propertyInput.UserID,
	}

	storage.DB.Create(&property)

	ctx.JSON(property)
}

func GetProperty(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	property := GetPropertyAndAssociationsByPropertyID(id, ctx)
	if property == nil {
		return
	}

	ctx.JSON(property)
}

func GetPropertiesByUserID(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	var properties []models.Property
	propertiesExist := storage.DB.Preload(clause.Associations).Where("user_id = ?", id).Find(&properties)

	if propertiesExist.Error != nil {
		utils.CreateError(
			iris.StatusInternalServerError,
			"Error", propertiesExist.Error.Error(), ctx)
		return
	}

	ctx.JSON(properties)
}

func DeleteProperty(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	propertyDeleted := storage.DB.Delete(&models.Property{}, id)

	if propertyDeleted.Error != nil {
		utils.CreateError(
			iris.StatusInternalServerError,
			"Error", propertyDeleted.Error.Error(), ctx)
		return
	}

	storage.DB.Where("property_id = ?", id).Delete(&models.Apartment{})
	ctx.StatusCode(iris.StatusNoContent)
}

func UpdateProperty(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	property := GetPropertyAndAssociationsByPropertyID(id, ctx)
	if property == nil {
		return
	}

	var propertyInput UpdatePropertyInput
	err := ctx.ReadJSON(&propertyInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var newApartments []models.Apartment
	var newApartmentImages []*[]string
	bedroomLow := property.BedroomLow
	bedroomHigh := property.BedroomHigh
	var bathroomLow float32 = property.BathroomLow
	var bathroomHigh float32 = property.BathroomHigh
	var rentLow float32 = propertyInput.Apartments[0].Rent
	var rentHigh float32 = propertyInput.Apartments[0].Rent

	for _, apartment := range propertyInput.Apartments {
		if apartment.Bathrooms < bathroomLow {
			bathroomLow = apartment.Bathrooms
		}
		if apartment.Bathrooms > bathroomHigh {
			bathroomHigh = apartment.Bathrooms
		}
		if *apartment.Bedrooms < bedroomLow {
			bedroomLow = *apartment.Bedrooms
		}
		if *apartment.Bedrooms > bedroomHigh {
			bedroomHigh = *apartment.Bedrooms
		}
		if apartment.Rent < rentLow {
			rentLow = apartment.Rent
		}
		if apartment.Rent > rentHigh {
			rentHigh = apartment.Rent
		}

		amenities, _ := json.Marshal(apartment.Amenities)

		currApartment := models.Apartment{
			Unit:        apartment.Unit,
			Bedrooms:    *apartment.Bedrooms,
			Bathrooms:   apartment.Bathrooms,
			PropertyID:  property.ID,
			SqFt:        apartment.SqFt,
			Rent:        apartment.Rent,
			Deposit:     *apartment.Deposit,
			LeaseLength: apartment.LeaseLength,
			AvailableOn: apartment.AvailableOn,
			Active:      apartment.Active,
			Amenities:   amenities,
			Description: apartment.Description,
		}

		if apartment.ID != nil {
			currApartment.ID = *apartment.ID
			updateApartmentAndImages(currApartment, apartment.Images)
		} else {
			newApartments = append(newApartments, currApartment)
			newApartmentImages = append(newApartmentImages, &apartment.Images)
		}
	}

	storage.DB.Create(&newApartments)

	for index, apartment := range newApartments {
		if len(*newApartmentImages[index]) > 0 {
			updateApartmentAndImages(apartment, *newApartmentImages[index])
		}
	}

	propertyAmenities, _ := json.Marshal(propertyInput.Amenities)
	includedUtilities, _ := json.Marshal(propertyInput.IncludedUtilities)

	property.UnitType = propertyInput.UnitType
	property.Description = propertyInput.Description
	property.IncludedUtilities = includedUtilities
	property.PetsAllowed = propertyInput.PetsAllowed
	property.LaundryType = propertyInput.LaundryType
	property.ParkingFee = *propertyInput.ParkingFee
	property.Amenities = propertyAmenities
	property.Name = propertyInput.Name
	property.FirstName = propertyInput.FirstName
	property.LastName = propertyInput.LastName
	property.Email = propertyInput.Email
	property.CallingCode = propertyInput.CallingCode
	property.CountryCode = propertyInput.CountryCode
	property.PhoneNumber = propertyInput.PhoneNumber
	property.Website = propertyInput.Website
	property.OnMarket = propertyInput.OnMarket
	property.BathroomHigh = bathroomHigh
	property.BathroomLow = bathroomLow
	property.BedroomLow = bedroomLow
	property.BedroomHigh = bedroomHigh
	property.RentLow = rentLow
	property.RentHigh = rentHigh

	imagesArr := insertImages(InsertImages{
		images:     propertyInput.Images,
		propertyID: strconv.FormatUint(uint64(property.ID), 10),
	})

	jsonImgs, _ := json.Marshal(imagesArr)

	property.Images = jsonImgs

	rowsUpdated := storage.DB.Model(&property).Updates(property)

	if rowsUpdated.Error != nil {
		utils.CreateError(
			iris.StatusInternalServerError,
			"Error", rowsUpdated.Error.Error(), ctx)
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

func GetPropertyAndAssociationsByPropertyID(id string, ctx iris.Context) *models.Property {

	var property models.Property
	propertyExists := storage.DB.Preload(clause.Associations).Find(&property, id)

	if propertyExists.Error != nil {
		utils.CreateInternalServerError(ctx)
		return nil
	}

	if propertyExists.RowsAffected == 0 {
		utils.CreateNotFound(ctx)
		return nil
	}

	return &property
}

func GetPropertiesByBoundingBox(ctx iris.Context) {
	var boundingBox BoundingBoxInput
	err := ctx.ReadJSON(&boundingBox)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var properties []models.Property
	storage.DB.Preload(clause.Associations).
		Where("lat >= ? AND lat <= ? AND lng >= ? AND lng <= ? AND on_market = true",
			boundingBox.LatLow, boundingBox.LatHigh, boundingBox.LngLow, boundingBox.LngHigh).
		Find(&properties)

	ctx.JSON(properties)
}

func updateApartmentAndImages(apartment models.Apartment, images []string) {
	apartmentID := strconv.FormatUint(uint64(apartment.ID), 10)

	apartmentImages := insertImages(InsertImages{
		images:      images,
		propertyID:  strconv.FormatUint(uint64(apartment.PropertyID), 10),
		apartmentID: &apartmentID,
	})

	if len(apartmentImages) > 0 {
		images, _ := json.Marshal(apartmentImages)
		apartment.Images = images
	}

	storage.DB.Model(&apartment).Updates(apartment)
}

func insertImages(arg InsertImages) []string {
	var imagesArr []string
	for _, image := range arg.images {
		if !strings.Contains(image, storage.BucketName) {
			imageID := randstr.Hex(16)
			imageStr := "property/" + arg.propertyID
			if arg.apartmentID != nil {
				imageStr += "/apartment/" + *arg.apartmentID
			}
			imageStr += "/" + imageID
			urlMap := storage.UploadBase64Image(image, imageStr)
			imagesArr = append(imagesArr, urlMap["url"])
		} else {
			imagesArr = append(imagesArr, image)
		}
	}
	return imagesArr
}

type InsertImages struct {
	images      []string
	propertyID  string
	apartmentID *string
}

type CreatePropertyInput struct {
	UnitType     string                 `json:"unitType" validate:"required,oneof=single multiple"`
	PropertyType string                 `json:"propertyType" validate:"required,max=256"`
	Street       string                 `json:"street" validate:"required,max=512"`
	City         string                 `json:"city" validate:"required,max=512"`
	State        string                 `json:"state" validate:"required,max=256"`
	Zip          int                    `json:"zip" validate:"required"`
	Lat          float32                `json:"lat" validate:"required"`
	Lng          float32                `json:"lng" validate:"required"`
	UserID       uint                   `json:"userID" validate:"required"`
	Apartments   []CreateApartmentInput `json:"apartments" validate:"required,dive"`
}

type CreateApartmentInput struct {
	Unit        string    `json:"unit" validate:"max=512"`
	Bedrooms    *int      `json:"bedrooms" validate:"gte=0,max=6,required"` // make int a pointer so 0 is accepted
	Bathrooms   float32   `json:"bathrooms" validate:"min=0.5,max=6.5,required"`
	Active      *bool     `json:"active" validate:"required"`
	AvailableOn time.Time `json:"availableOn" validate:"required"`
}

type UpdatePropertyInput struct {
	UnitType          string                  `json:"unitType" validate:"required,oneof=single multiple"`
	Description       string                  `json:"description"`
	Images            []string                `json:"images"`
	IncludedUtilities []string                `json:"includedUtilities"`
	PetsAllowed       string                  `json:"petsAllowed" validate:"required"`
	LaundryType       string                  `json:"laundryType" validate:"required"`
	ParkingFee        *float32                `json:"parkingFee"`
	Amenities         []string                `json:"amenities"`
	Name              string                  `json:"name"`
	FirstName         string                  `json:"firstName"`
	LastName          string                  `json:"lastName"`
	Email             string                  `json:"email" validate:"required,email"`
	CallingCode       string                  `json:"callingCode"`
	CountryCode       string                  `json:"countryCode"`
	PhoneNumber       string                  `json:"phoneNumber" validate:"required"`
	Website           string                  `json:"website" validate:"omitempty,url"`
	OnMarket          *bool                   `json:"onMarket" validate:"required"`
	Apartments        []UpdateApartmentsInput `json:"apartments" validate:"required,dive"`
}

type UpdateApartmentsInput struct {
	ID          *uint     `json:"ID"`
	Unit        string    `json:"unit" validate:"max=512"`
	Bedrooms    *int      `json:"bedrooms" validate:"gte=0,max=6,required"` // make int a pointer so 0 is accepted
	Bathrooms   float32   `json:"bathrooms" validate:"min=0.5,max=6.5,required"`
	SqFt        int       `json:"sqFt" validate:"max=100000000000,required"`
	Rent        float32   `json:"rent" validate:"required"`
	Deposit     *float32  `json:"deposit" validate:"required"`
	LeaseLength string    `json:"leaseLength" validate:"required,max=256"`
	AvailableOn time.Time `json:"availableOn" validate:"required"`
	Active      *bool     `json:"active" validate:"required"`
	Images      []string  `json:"images"`
	Amenities   []string  `json:"amenities"`
	Description string    `json:"description"`
}

type BoundingBoxInput struct {
	LatLow  float32 `json:"latLow" validate:"required"`
	LatHigh float32 `json:"latHigh" validate:"required"`
	LngLow  float32 `json:"lngLow" validate:"required"`
	LngHigh float32 `json:"lngHigh" validate:"required"`
}
