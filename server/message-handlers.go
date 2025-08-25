package server

import (
	"fluffy-coto-tribble/server/services"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func CreateMessage(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg services.Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		id := fmt.Sprintf("m_%s", ShortUUID())
		now := time.Now().Unix()

		newMessage := services.Message{
			ID:        id,
			ChatID:    msg.ChatID,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			Media:     msg.Media,
			Timestamp: now,
		}

		if err := services.CreateMessage(client, "messages", newMessage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Message created successfully",
			"data":    newMessage,
		})
	}
}

func GetChatMessage(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("chatId")
		msgID := c.Param("id")

		msg, err := services.GetChatMessage(client, "messages", chatID, msgID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}

func GetAllChatMessages(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("chatId")

		messages, err := services.GetAllChatMessages(client, "messages", chatID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

func UpdateMessage(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("chatId")
		msgID := c.Param("id")

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		avUpdates, err := attributevalue.MarshalMap(updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updates"})
			return
		}

		if err := services.UpdateMessage(client, "messages", chatID, msgID, avUpdates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
	}
}

func DeleteMessage(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("chatId")
		msgID := c.Param("id")

		if err := services.DeleteMessage(client, "messages", chatID, msgID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
	}
}
