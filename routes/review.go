package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"
	"math"
	"strconv"

	"github.com/kataras/iris/v12"
)

func CreateReview(ctx iris.Context) {
	params := ctx.Params()
	propertyID := params.Get("id")

	property := GetPropertyAndAssociationsByPropertyID(propertyID, ctx)

	if property == nil {
		return
	}

	var reviewInput CreateReviewInput
	err := ctx.ReadJSON(&reviewInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	propID, convErr := strconv.ParseUint(propertyID, 10, 32)

	if convErr != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	review := models.Review{
		UserID:     reviewInput.UserID,
		PropertyID: uint(propID),
		Title:      reviewInput.Title,
		Body:       reviewInput.Body,
		Stars:      reviewInput.Stars,
	}

	storage.DB.Create(&review)

	updatePropertyStars(property, float32(review.Stars))

	ctx.JSON(review)
}

func updatePropertyStars(property *models.Property, stars float32) {
	var avg float32
	reviewsLength := len(property.Reviews)
	if reviewsLength == 0 {
		avg = stars
	} else {
		var sum float32
		for i := 0; i < len(property.Reviews); i++ {
			sum += float32(property.Reviews[i].Stars)
		}

		avg = (sum + stars) / (float32(reviewsLength) + 1)
	}

	avg = float32(math.Round(float64(avg)*10) / 10)
	storage.DB.Model(&property).Update("stars", avg)
}

type CreateReviewInput struct {
	UserID uint   `json:"userID" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
	Stars  int    `json:"stars" validate:"required,gt=0,lt=6"`
}
