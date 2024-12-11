package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func main() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	collection = client.Database("naruto_db").Collection("characters")

	router := gin.Default()

	router.GET("/character", indexUser)
	router.POST("/character", createUser)
	router.GET("/character/:slug", readUser)
	router.PUT("/character/:slug", updateUser)
	router.DELETE("/character/:slug", deleteUser)

	router.Run(":8001")
}

func indexUser(c *gin.Context) {
	var users []bson.M

	// Ambil semua data tanpa pagination
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	// Menyimpan hasil query ke dalam slice `users`
	for cursor.Next(context.TODO()) {
		var user bson.M
		err := cursor.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	// Kembalikan hasil dalam response
	c.JSON(http.StatusOK, gin.H{
		"messages": "Success retrieved all data",
		"result":   users,
	})
}

func createUser(c *gin.Context) {
	// Definisikan tipe data secara langsung di dalam fungsi
	var user struct {
		Name     string   `json:"name"`
		Images   []string `json:"images"`
		Personal struct {
			Birthdate   string `json:"birthdate"`
			Sex         string `json:"sex"`
			Status      string `json:"status"`
			Height      string `json:"height"`
			Weight      string `json:"weight"`
			BloodType   string `json:"bloodType"`
			Occupation  string `json:"occupation"`
			Affiliation string `json:"affiliation"`
			Clan        string `json:"clan"` // Menambahkan clan
		} `json:"personal"`
		Rank struct {
			NinjaRank string `json:"ninjaRank"`
		} `json:"rank"`
		Debut struct {
			Anime     string `json:"anime"`
			AppearsIn string `json:"appearsIn"`
		} `json:"debut"` // Menambahkan debut
		Jutsu []string `json:"jutsu"` // Menambahkan jutsu
	}

	// Bind data dari JSON request body ke variabel user
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug berdasarkan nama
	userSlug := slug.Make(user.Name)

	// Insert data user ke MongoDB
	result, err := collection.InsertOne(context.TODO(), bson.M{
		"name":     user.Name,
		"slug":     userSlug,
		"images":   user.Images,
		"personal": user.Personal,
		"rank":     user.Rank,
		"debut":    user.Debut,
		"jutsu":    user.Jutsu,
	})

	// Handle error jika insert gagal
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan response dengan hasil insert
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func readUser(c *gin.Context) {
	slugParam := c.Param("slug")
	var user bson.M

	err := collection.FindOne(context.TODO(), bson.M{"slug": slugParam}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": "Success retrieved data",
		"result":   user,
	})
}

func updateUser(c *gin.Context) {
	// Ambil slug dari URL parameter
	slugParam := c.Param("slug")

	// Definisikan tipe data untuk update user
	var user struct {
		Name     string   `json:"name"`
		Images   []string `json:"images"`
		Personal struct {
			Birthdate   string `json:"birthdate"`
			Sex         string `json:"sex"`
			Status      string `json:"status"`
			Height      string `json:"height"`
			Weight      string `json:"weight"`
			BloodType   string `json:"bloodType"`
			Occupation  string `json:"occupation"`
			Affiliation string `json:"affiliation"`
			Clan        string `json:"clan"` // Menambahkan clan
		} `json:"personal"`
		Rank struct {
			NinjaRank string `json:"ninjaRank"`
		} `json:"rank"`
		Debut struct {
			Anime     string `json:"anime"`
			AppearsIn string `json:"appearsIn"`
		} `json:"debut"` // Menambahkan debut
		Jutsu []string `json:"jutsu"` // Menambahkan jutsu
	}

	// Bind data dari JSON request body ke variabel user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate slug baru berdasarkan nama yang diperbarui
	newSlug := slug.Make(user.Name)

	// Filter untuk mencari user berdasarkan slug yang lama
	filter := bson.M{"slug": slugParam}

	// Update data user di MongoDB dengan data yang baru
	update := bson.M{
		"$set": bson.M{
			"name":     user.Name,
			"slug":     newSlug, // Mengganti slug dengan yang baru
			"images":   user.Images,
			"personal": user.Personal,
			"rank":     user.Rank,
			"debut":    user.Debut,
			"jutsu":    user.Jutsu,
		},
	}

	// Melakukan update pada database
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Mencari user yang telah diperbarui
	var updatedUser struct {
		Name     string   `json:"name"`
		Images   []string `json:"images"`
		Personal struct {
			Birthdate   string `json:"birthdate"`
			Sex         string `json:"sex"`
			Status      string `json:"status"`
			Height      string `json:"height"`
			Weight      string `json:"weight"`
			BloodType   string `json:"bloodType"`
			Occupation  string `json:"occupation"`
			Affiliation string `json:"affiliation"`
			Clan        string `json:"clan"`
		} `json:"personal"`
		Rank struct {
			NinjaRank string `json:"ninjaRank"`
		} `json:"rank"`
		Debut struct {
			Anime     string `json:"anime"`
			AppearsIn string `json:"appearsIn"`
		} `json:"debut"`
		Jutsu []string `json:"jutsu"`
		Slug  string   `json:"slug"`
	}

	// Retrieve updated user
	err = collection.FindOne(context.TODO(), filter).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated user"})
		return
	}

	// Response setelah update berhasil
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"messages": "User updated",
		"result":   updatedUser,
	})
}

func deleteUser(c *gin.Context) {
	slugParam := c.Param("slug")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"slug": slugParam})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Character deleted"})
}
