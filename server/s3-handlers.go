package server

import (
	"fluffy-coto-tribble/server/authentication"
	"fluffy-coto-tribble/server/services"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func Upload(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		claims := authentication.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		// Parse file
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
			return
		}
		defer file.Close()

		// Upload file
		fileURL, err := services.UploadFile(client, header.Filename, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"url":     fileURL,
		})
	}
}

func Download(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			return
		}

		// Auth
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		claims := authentication.ParseAccessToken(token)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		filename := c.Query("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filename parameter"})
			return
		}

		url, err := services.DownloadFile(client, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"downloadUrl": url,
		})
	}
}
