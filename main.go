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

	r := gin.Default()

	r.GET("/character", character.IndexUser)
	r.GET("/character/search", character.SearchCharacter)
	r.POST("/character", character.CreateUser)
	r.GET("/character/:slug", character.ReadUser)
	r.PUT("/character/:slug", character.UpdateUser)
	r.DELETE("/character/:slug", character.DeleteUser)

	r.GET("/tailedbeast", tailedbeast.IndexTailedBeast)
	r.GET("/tailedbeast/search", tailedbeast.SearchTailedBeast)
	r.POST("/tailedbeast", tailedbeast.CreateTailedBeast)
	r.GET("/tailedbeast/:slug", tailedbeast.ReadTailedBeast)
	r.PUT("/tailedbeast/:slug", tailedbeast.UpdateTailedBeast)
	r.DELETE("/tailedbeast/:slug", tailedbeast.DeleteTailedBeast)

	r.Run(":8001")

}
