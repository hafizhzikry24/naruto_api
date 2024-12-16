package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my-gin-app/character"
	"my-gin-app/tailedbeast"
)

var client *mongo.Client
var db *mongo.Database
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(os.Getenv("MONGO_DB"))
	character.InitCharacterCollection(db.Collection(os.Getenv("MONGO_COLLECTION")))
	tailedbeast.InitTailedBeastCollection(db.Collection(os.Getenv("MONGO_COLLECTION_TAILEDBEAST")))

	router := gin.Default()
	protected := router.Group("/")
	protected.Use(JWTAuthMiddleware)

	router.POST("/login", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if creds.Username != "admin" || creds.Password != "password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})

	})

	protected.GET("/character", character.IndexUser)
	protected.GET("/character/search", character.SearchCharacter)
	protected.POST("/character", character.CreateUser)
	protected.GET("/character/:slug", character.ReadUser)
	protected.PUT("/character/:slug", character.UpdateUser)
	protected.DELETE("/character/:slug", character.DeleteUser)

	protected.GET("/tailedbeast", tailedbeast.IndexTailedBeast)
	protected.GET("/tailedbeast/search", tailedbeast.SearchTailedBeast)
	protected.POST("/tailedbeast", tailedbeast.CreateTailedBeast)
	protected.GET("/tailedbeast/:slug", tailedbeast.ReadTailedBeast)
	protected.PUT("/tailedbeast/:slug", tailedbeast.UpdateTailedBeast)
	protected.DELETE("/tailedbeast/:slug", tailedbeast.DeleteTailedBeast)

	router.Run(":8001")

}

func JWTAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
		c.Abort()
		return
	}

	c.Set("user", claims.Username)

	c.Next()
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
