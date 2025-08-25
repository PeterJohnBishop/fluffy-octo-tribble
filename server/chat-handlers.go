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

func CreateChat(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var chat services.Chat
		if err := c.ShouldBindJSON(&chat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		id := fmt.Sprintf("c_%s", ShortUUID())
		now := time.Now().Unix()

		newChat := services.Chat{
			ID:          id,
			Users:       chat.Users,
			Messages:    []string{}, // start empty
			DateCreated: now,
			DateUpdated: now,
		}

		if err := services.CreateChat(client, "chats", newChat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Chat created successfully",
			"chat":    newChat,
		})
	}
}

func GetAllChats(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chats, err := services.GetAllChats(client, "chats")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"chats": chats})
	}
}

func GetChatById(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("id")

		chat, err := services.GetChatById(client, "chats", chatID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"chat": chat})
	}
}

func UpdateChat(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("id")
		var updates map[string]interface{}

		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		updates["dateUpdated"] = time.Now().Unix()

		avUpdates, err := attributevalue.MarshalMap(updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updates"})
			return
		}

		if err := services.UpdateChat(client, "chats", chatID, avUpdates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Chat updated successfully"})
	}
}

func DeleteChat(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.Param("id")

		if err := services.DeleteChat(client, "chats", chatID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
	}
}
