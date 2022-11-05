package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func CreateMessage(ctx iris.Context) {
	var req CreateMessageInput

	err := ctx.ReadJSON(&req)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	claims := jwt.Get(ctx).(*utils.AccessToken)

	if req.SenderID != claims.ID {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	message := models.Message{
		ConversationID: req.ConversationID,
		SenderID:       req.SenderID,
		ReceiverID:     req.ReceiverID,
		Text:           req.Text,
	}

	storage.DB.Create(&message)

	ctx.JSON(message)
}

type CreateMessageInput struct {
	ConversationID uint   `json:"conversationID" validate:"required"`
	SenderID       uint   `json:"senderID" validate:"required"`
	ReceiverID     uint   `json:"receiverID" validate:"required"`
	Text           string `json:"text" validate:"required,lt=5000"`
}
