package server

import (
	"fluffy-coto-tribble/server/authentication"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
)

func AddDynamoDBRoutes(client *dynamodb.Client, r *gin.Engine) {
	r.POST("/register", CreateUser(client))
	r.POST("/login", AuthUser(client))
	r.POST("/refresh-token", authentication.RefreshTokenHandler(client))

	auth := r.Group("/", authentication.AuthMiddleware())
	{
		// users
		auth.GET("/users", GetAllUsers(client))
		auth.GET("/users/:id", GetUserByID(client))
		auth.PUT("/users", UpdateUser(client))
		auth.PUT("/users/password", UpdatePassword(client))
		auth.DELETE("/users/:id", DeleteUser(client))
		// chats
		auth.POST("/chats", CreateChat(client))
		auth.GET("/chats", GetAllChats(client))
		auth.GET("/chats/:id", GetChatById(client))
		auth.PUT("/chats/:id", UpdateChat(client))
		auth.DELETE("/chats/:id", DeleteChat(client))
		// messages
		auth.POST("/messages", CreateMessage(client))
		auth.GET("/messages/:chatId/:id", GetChatMessage(client))
		auth.GET("/messages/:chatId", GetAllChatMessages(client))
		auth.PUT("/messages/:chatId/:id", UpdateMessage(client))
		auth.DELETE("/messages/:chatId/:id", DeleteMessage(client))
	}
}

func AddMapRoutes(client *maps.Client, r *gin.Engine) {
	auth := r.Group("/", authentication.AuthMiddleware())
	{
		auth.GET("/geocode", Geocode(client))
		auth.GET("/reverse-geocode", ReverseGeocode(client))
		auth.GET("/directions", GetDirections(client))
	}
}

func AddS3Routes(s3client *s3.Client, dynamoclient *dynamodb.Client, r *gin.Engine) {
	auth := r.Group("/", authentication.AuthMiddleware())
	{
		auth.POST("/upload", Upload(s3client, dynamoclient))
		auth.GET("/files", GetUserFilesHandler(dynamoclient, s3.NewPresignClient(s3client)))
		auth.GET("/download", Download(s3client))
	}
}
