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

	router.GET("/character", character.IndexUser)
	router.GET("/character/search", character.SearchCharacter)
	router.POST("/character", character.CreateUser)
	router.GET("/character/:slug", character.ReadUser)
	router.PUT("/character/:slug", character.UpdateUser)
	router.DELETE("/character/:slug", character.DeleteUser)

	router.GET("/tailedbeast", tailedbeast.IndexTailedBeast)
	router.GET("/tailedbeast/search", tailedbeast.SearchTailedBeast)
	router.POST("/tailedbeast", tailedbeast.CreateTailedBeast)
	router.GET("/tailedbeast/:slug", tailedbeast.ReadTailedBeast)
	router.PUT("/tailedbeast/:slug", tailedbeast.UpdateTailedBeast)
	router.DELETE("/tailedbeast/:slug", tailedbeast.DeleteTailedBeast)

	router.Run(":8001")

}
