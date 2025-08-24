package server

import (
	"fluffy-coto-tribble/server/services"
	"log"

	"github.com/gin-gonic/gin"
)

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// // connect with DynamoDB
	dynamoClient := services.ConnectDB()
	AddDynamoDBRoutes(dynamoClient, router)

	// // connect with S3
	s3Client := services.ConnectS3()
	AddS3Routes(s3Client, router)

	// // connect with Google Maps
	mapClient := services.FindMaps()
	AddMapRoutes(mapClient, router)

	log.Println("Starting listening on :8080")
	router.Run(":8080")
}
