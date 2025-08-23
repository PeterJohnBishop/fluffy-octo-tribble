package server

import (
	"fluffy-coto-tribble/server/services"
	"log"

	"github.com/gin-gonic/gin"
)

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	cfg := services.StartAws()

	// // connect with DynamoDB
	dynamoClient := services.ConnectDB(cfg)
	AddDynamoDBRoutes(dynamoClient, router)

	// // connect with S3
	// s3Client := services.ConnectS3(cfg)

	// // connect with Google Maps
	// mapClient := services.FindMaps()

	log.Println("Starting listening on :8080")
	router.Run(":8080")
}
