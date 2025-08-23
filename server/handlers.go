package server

import (
	"fluffy-coto-tribble/server/authentication"
	"fluffy-coto-tribble/server/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

func CreateUser(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user services.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		id, err := uuid.NewV1()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating user id"})
			return
		}

		email := strings.ToLower(user.Email)
		userId := fmt.Sprintf("u_%s", id)

		hashedPassword, err := authentication.HashedPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		newUser := map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: userId},
			"name":     &types.AttributeValueMemberS{Value: user.Name},
			"email":    &types.AttributeValueMemberS{Value: email},
			"password": &types.AttributeValueMemberS{Value: hashedPassword},
		}

		if err := services.CreateUser(client, "users", newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user.id": userId,
		})
	}
}

func AuthUser(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := services.GetUserByEmail(client, "users", req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No user found with that email."})
			return
		}

		if !authentication.CheckPasswordHash(req.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password verification failed"})
			return
		}

		userClaims := authentication.UserClaims{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			},
		}

		token, err := authentication.NewAccessToken(userClaims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err := authentication.NewRefreshToken(userClaims.StandardClaims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Login Success",
			"token":         token,
			"refresh_token": refreshToken,
			"user":          user,
		})
	}
}

func GetAllUsers(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}
		if authentication.ParseAccessToken(token) == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		resp, err := services.GetAllUsers(client, "users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
			return
		}

		var users []services.User
		for _, item := range resp {
			var user services.User
			if err := attributevalue.UnmarshalMap(item, &user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Users Found!",
			"users":   users,
		})
	}
}

func GetUserByID(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" || authentication.ParseAccessToken(token) == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		resp, err := services.GetUserById(client, "users", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		var user services.User
		if err := attributevalue.UnmarshalMap(resp, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User Found!",
			"user":    user,
		})
	}
}

func UpdateUser(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" || authentication.ParseAccessToken(token) == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		var user services.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if err := services.UpdateUser(client, "users", user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User Updated!"})
	}
}

func UpdatePassword(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" || authentication.ParseAccessToken(token) == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		var user services.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		hashedPassword, err := authentication.HashedPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword

		if err := services.UpdatePassword(client, "users", user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User Password Updated!"})
	}
}

func DeleteUser(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" || authentication.ParseAccessToken(token) == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token"})
			return
		}

		if err := services.DeleteUser(client, "users", id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User Deleted!"})
	}
}
