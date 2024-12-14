package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my-gin-app/character"
	"my-gin-app/tailedbeast"
)

var client *mongo.Client
var db *mongo.Database

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
	protected.Use(ApiKeyMiddleware)

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

func ApiKeyMiddleware(c *gin.Context) {
	apiKeyHeader := c.GetHeader("X-API-KEY")

	if apiKeyHeader != os.Getenv("X_API_KEY") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
