package routes

import (
	"apartments-clone-server/models"
	"apartments-clone-server/storage"
	"apartments-clone-server/utils"
	"errors"
	"sort"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func CreateConversation(ctx iris.Context) {
	var req CreateConversationInput

	err := ctx.ReadJSON(&req)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var prevConversation models.Conversation
	conversationExists := storage.DB.
		Where("property_id = ? AND tenant_id = ? AND owner_id = ? ", req.PropertyID, req.TenantID, req.OwnerID).
		Find(&prevConversation)

	if conversationExists.Error != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	if conversationExists.RowsAffected > 0 {
		ctx.StatusCode(iris.StatusConflict)
		ctx.Text("Conversation already exists")
		return
	}

	var messages []models.Message
	messages = append(messages, models.Message{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Text:       req.Text,
	})

	conversation := models.Conversation{
		TenantID:   req.TenantID,
		OwnerID:    req.OwnerID,
		PropertyID: req.PropertyID,
		Messages:   messages,
	}

	storage.DB.Create(&conversation)

	ctx.JSON(conversation)
}

func GetConversationByID(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	result, err := getConversationResult(id, ctx)

	if err != nil {
		return
	}

	var messages []models.Message
	messagesQuery := storage.DB.Where("conversation_id = ?", id).Order("created_at DESC").Find(&messages)

	if messagesQuery.Error != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	result.Messages = messages

	ctx.JSON(result)
}

func GetConversationsByUserID(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	results, err := getConversationResultsByUserID(id, ctx)

	if err != nil {
		return
	}

	var conversationIDs []uint
	for _, conversation := range results {
		conversationIDs = append(conversationIDs, conversation.ID)
	}

	var messages []models.Message

	messagesQuery := storage.DB.Raw(`
		SELECT messages.* 
		FROM messages
		INNER JOIN (
			SELECT conversation_id, MAX(created_at) AS created_at
			FROM messages
			WHERE conversation_id IN ? 
			GROUP BY conversation_id
		) AS recentMessages
		ON messages.conversation_id = recentMessages.conversation_id 
		AND messages.created_at = recentMessages.created_at`, conversationIDs).Scan(&messages)

	messageMap := make(map[uint][]models.Message)
	for _, message := range messages {
		messageSlice := []models.Message{message}
		messageMap[message.ConversationID] = messageSlice
	}

	for index, conversation := range results {
		results[index].Messages = messageMap[conversation.ID]
	}

	if messagesQuery.Error != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	sort.Slice(results, func(i int, j int) bool {
		return results[i].Messages[0].CreatedAt.After(results[j].Messages[0].CreatedAt)
	})

	ctx.JSON(results)
}

func getConversationResult(id string, ctx iris.Context) (ConversationResult, error) {
	var result ConversationResult
	resultQuery := storage.DB.Table("conversations").
		Select(`conversations.*,
		 properties.name, properties.street, properties.city, properties.state, 
		 owners.first_name as owner_first_name, owners.last_name as owner_last_name, owners.email as owner_email,
		 tenants.first_name as tenant_first_name, tenants.last_name as tenant_last_name, tenants.email as tenant_email`).
		Joins("INNER JOIN properties on properties.id = conversations.property_id").
		Joins("INNER JOIN users owners on conversations.owner_id = owners.id").
		Joins("INNER JOIN users tenants on conversations.tenant_id = tenants.id").
		Where("conversations.id = ?", id).
		Scan(&result)

	if resultQuery.Error != nil {
		utils.CreateInternalServerError(ctx)
		return result, resultQuery.Error
	}

	if resultQuery.RowsAffected == 0 {
		utils.CreateNotFound(ctx)
		return result, errors.New("Result not found")
	}

	return result, nil
}

func getConversationResultsByUserID(id string, ctx iris.Context) ([]ConversationResult, error) {
	var result []ConversationResult
	resultQuery := storage.DB.Table("conversations").
		Select(`conversations.*,
		 properties.name, properties.street, properties.city, properties.state, 
		 owners.first_name as owner_first_name, owners.last_name as owner_last_name, owners.email as owner_email,
		 tenants.first_name as tenant_first_name, tenants.last_name as tenant_last_name, tenants.email as tenant_email`).
		Joins("INNER JOIN properties on properties.id = conversations.property_id").
		Joins("INNER JOIN users owners on conversations.owner_id = owners.id").
		Joins("INNER JOIN users tenants on conversations.tenant_id = tenants.id").
		Where("conversations.tenant_id = ?", id).Or("conversations.owner_id = ?", id).
		Scan(&result)

	if resultQuery.Error != nil {
		utils.CreateInternalServerError(ctx)
		return result, resultQuery.Error
	}

	if resultQuery.RowsAffected == 0 {
		utils.CreateNotFound(ctx)
		return result, errors.New("Result not found")
	}

	return result, nil
}

type ConversationResult struct {
	// Conversation
	gorm.Model
	TenantID   uint `json:"tenantID"`
	OwnerID    uint `json:"ownerID"`
	PropertyID uint `json:"propertyID"`
	// Property
	Name   string `json:"propertyName"`
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	// Owner / User
	OwnerFirstName string `json:"ownerFirstName"`
	OwnerLastName  string `json:"ownerLastName"`
	OwnerEmail     string `json:"ownerEmail"`
	// Tenenant / User
	TenantFirstName string `json:"tenantFirstName"`
	TenantLastName  string `json:"tenantLastName"`
	TenantEmail     string `json:"tenantEmail"`
	// Conversation / Message
	Messages []models.Message `gorm:"foreignKey:ID" json:"messages"`
}

type CreateConversationInput struct {
	TenantID   uint   `json:"tenantID" validate:"required"`
	OwnerID    uint   `json:"ownerID" validate:"required"`
	PropertyID uint   `json:"propertyID" validate:"required"`
	SenderID   uint   `json:"senderID" validate:"required"`
	ReceiverID uint   `json:"receiverID" validate:"required"`
	Text       string `json:"text" validate:"required,lt=5000"`
}
