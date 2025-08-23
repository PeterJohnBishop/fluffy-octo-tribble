package server

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func AddDynamoDBRoutes(client *dynamodb.Client, r *gin.Engine) {
	r.POST("/register", CreateUser(client))
	r.POST("/login", AuthUser(client))
	r.GET("/users", GetAllUsers(client))
	r.GET("/users/:id", GetUserByID(client))
	r.PUT("/users", UpdateUser(client))
	r.PUT("/users/password", UpdatePassword(client))
	r.DELETE("/users/:id", DeleteUser(client))
}
