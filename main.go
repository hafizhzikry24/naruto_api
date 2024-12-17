package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my-gin-app/character"
	"my-gin-app/middleware"
	"my-gin-app/tailedbeast"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database(os.Getenv("MONGO_DB"))

	characterRepo := character.NewRepository(db.Collection(os.Getenv("MONGO_COLLECTION")))
	tailedBeastRepo := tailedbeast.NewRepository(db.Collection(os.Getenv("MONGO_COLLECTION_TAILEDBEAST")))

	characterService := character.NewService(characterRepo)
	tailedBeastService := tailedbeast.NewService(tailedBeastRepo)

	characterHandler := character.NewHandler(characterService)
	tailedBeastHandler := tailedbeast.NewHandler(tailedBeastService)

	router := gin.Default()
	router.Use(middleware.APIKeyMiddleware())

	router.GET("/character", characterHandler.IndexUser)
	router.GET("/character/search", characterHandler.SearchCharacter)
	router.POST("/character", characterHandler.CreateUser)
	router.GET("/character/:slug", characterHandler.ReadUser)
	router.PUT("/character/:slug", characterHandler.UpdateUser)
	router.DELETE("/character/:slug", characterHandler.DeleteUser)

	router.GET("/tailedbeast", tailedBeastHandler.IndexTailedBeast)
	router.GET("/tailedbeast/search", tailedBeastHandler.SearchTailedBeast)
	router.POST("/tailedbeast", tailedBeastHandler.CreateTailedBeast)
	router.GET("/tailedbeast/:slug", tailedBeastHandler.ReadTailedBeast)
	router.PUT("/tailedbeast/:slug", tailedBeastHandler.UpdateTailedBeast)
	router.DELETE("/tailedbeast/:slug", tailedBeastHandler.DeleteTailedBeast)

	router.Run(":8001")
	if err := router.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
