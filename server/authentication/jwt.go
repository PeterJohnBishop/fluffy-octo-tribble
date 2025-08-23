package authentication

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), error
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     = time.Minute * 15
	RefreshTokenTTL    = time.Hour * 24 * 7
)

// Load .env once at startup
func InitAuth() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	envPath := filepath.Join(".", ".env") // this points to fluffy-octo-tribble/.env

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file at %s: %v", envPath, err)
	}

	AccessTokenSecret = os.Getenv("TOKEN_SECRET")
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")

	if AccessTokenSecret == "" || RefreshTokenSecret == "" {
		log.Fatal("TOKEN_SECRET or REFRESH_TOKEN_SECRET is missing")
	}
}

type UserClaims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(AccessTokenSecret))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(RefreshTokenSecret))
}

func ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AccessTokenSecret), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		fmt.Println("Token verification failed:", err) // Debugging output
		return nil
	}

	claims, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok {
		fmt.Println("Failed to cast token claims")
		return nil
	}

	return claims
}

func ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(RefreshTokenSecret), nil
	})
	if err != nil || !parsedRefreshToken.Valid {
		fmt.Println("Refresh token verification failed:", err)
		return nil
	}

	claims, ok := parsedRefreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		fmt.Println("Failed to cast refresh token claims")
		return nil
	}

	return claims
}

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow /register and /login without auth
		if strings.HasPrefix(c.Request.URL.Path, "/register") ||
			strings.HasPrefix(c.Request.URL.Path, "/login") {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication Header is missing!"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userClaims := ParseAccessToken(token)
		if userClaims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token!"})
			return
		}

		// You can attach claims to context for use in handlers
		c.Set("userClaims", userClaims)

		c.Next()
	}
}

type VerifyRefreshRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func VerifyRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyRefreshRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		claims := ParseRefreshToken(req.Token)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token!"})
			return
		}

		// Save user ID in Gin context
		c.Set("userID", req.ID)

		c.Next()
	}
}
