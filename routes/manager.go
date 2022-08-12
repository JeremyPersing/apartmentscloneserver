package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
)

func CreateManager(ctx iris.Context) {
	const maxSize = 10 * iris.MB
	ctx.SetMaxRequestBodySize(maxSize)

	var managerInput ManagerInput
	err := ctx.ReadJSON(&managerInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var url string = ""
	if managerInput.Image != "" {
		res := storage.UploadBase64Image(
			managerInput.Image,
			strings.ReplaceAll(fmt.Sprint(managerInput.UserID)+"/"+managerInput.Name, " ", ""),
		)

		url = res["url"]
	}

	manager := models.Manager{
		Name:        managerInput.Name,
		UserID:      managerInput.UserID,
		Email:       managerInput.Email,
		PhoneNumber: managerInput.PhoneNumber,
		Website:     managerInput.Website,
		Image:       url,
	}
	storage.DB.Create(&manager)

	ctx.JSON(manager)
}

func GetManagerByUserID(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	var manager models.Manager
	managerExists := storage.DB.Where("user_id = ?", id).Find(&manager)

	if managerExists.RowsAffected == 0 {
		utils.CreateError(iris.StatusNotFound, "Manager Not Found", "Manager Not Found", ctx)
		return
	}

	ctx.JSON(manager)
}

type ManagerInput struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phoneNumber"`
	Website     string `json:"website"`
	Image       string `json:"image"`
	UserID      uint   `json:"userId" validate:"required"`
}
