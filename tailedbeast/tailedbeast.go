package tailedbeast

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tailedbeastCollection *mongo.Collection

func InitTailedBeastCollection(coll *mongo.Collection) {
	tailedbeastCollection = coll
}

func IndexTailedBeast(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	var tailedBeasts []bson.M

	if page == "" && limit == "" {
		cursor, err := tailedbeastCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var tailedBeast bson.M
			err := cursor.Decode(&tailedBeast)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tailedBeasts = append(tailedBeasts, tailedBeast)
		}
	} else {
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
			return
		}

		skip := (pageInt - 1) * limitInt

		cursor, err := tailedbeastCollection.Find(context.TODO(), bson.D{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limitInt)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var tailedBeast bson.M
			err := cursor.Decode(&tailedBeast)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tailedBeasts = append(tailedBeasts, tailedBeast)
		}

		count, err := tailedbeastCollection.CountDocuments(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalPages := (count + int64(limitInt) - 1) / int64(limitInt)

		c.JSON(http.StatusOK, gin.H{
			"messages":   "Success retrieved data",
			"result":     tailedBeasts,
			"page":       pageInt,
			"limit":      limitInt,
			"totalPages": totalPages,
			"totalItems": count,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": "Success retrieved all data",
		"result":   tailedBeasts,
	})
}

// searchTailedBeast: Menyaring tailedbeast berdasarkan nama
func SearchTailedBeast(c *gin.Context) {
	nameQuery := c.DefaultQuery("name", "")
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
		return
	}

	filter := bson.M{
		"name": bson.M{
			"$regex":   nameQuery,
			"$options": "i", // Case insensitive
		},
	}

	var tailedBeasts []bson.M
	cursor, err := tailedbeastCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var tailedBeast bson.M
		if err := cursor.Decode(&tailedBeast); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tailedBeasts = append(tailedBeasts, tailedBeast)
	}

	if len(tailedBeasts) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"messages": "Found tailed beasts",
			"result":   tailedBeasts,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "No tailed beasts found"})
	}
}

// createTailedBeast: Menambahkan tailedbeast baru
func CreateTailedBeast(c *gin.Context) {
	var tailedBeast struct {
		Name        string   `json:"name"`
		Images      []string `json:"images"`
		Rank        string   `json:"rank"`
		Abilities   []string `json:"abilities"`
		Personality string   `json:"personality"`
	}

	if err := c.BindJSON(&tailedBeast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tailedBeastSlug := slug.Make(tailedBeast.Name)

	result, err := tailedbeastCollection.InsertOne(context.TODO(), bson.M{
		"name":        tailedBeast.Name,
		"slug":        tailedBeastSlug,
		"images":      tailedBeast.Images,
		"rank":        tailedBeast.Rank,
		"abilities":   tailedBeast.Abilities,
		"personality": tailedBeast.Personality,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": result})
}

// readTailedBeast: Menampilkan tailedbeast berdasarkan slug
func ReadTailedBeast(c *gin.Context) {
	slugParam := c.Param("slug")
	var tailedBeast bson.M

	err := tailedbeastCollection.FindOne(context.TODO(), bson.M{"slug": slugParam}).Decode(&tailedBeast)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tailed Beast Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": "Success retrieved data",
		"result":   tailedBeast,
	})
}

// updateTailedBeast: Memperbarui tailedbeast
func UpdateTailedBeast(c *gin.Context) {
	slugParam := c.Param("slug")

	var tailedBeast struct {
		Name        string   `json:"name"`
		Images      []string `json:"images"`
		Rank        string   `json:"rank"`
		Abilities   []string `json:"abilities"`
		Personality string   `json:"personality"`
	}

	if err := c.ShouldBindJSON(&tailedBeast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSlug := slug.Make(tailedBeast.Name)

	filter := bson.M{"slug": slugParam}
	update := bson.M{
		"$set": bson.M{
			"name":        tailedBeast.Name,
			"slug":        newSlug,
			"images":      tailedBeast.Images,
			"rank":        tailedBeast.Rank,
			"abilities":   tailedBeast.Abilities,
			"personality": tailedBeast.Personality,
		},
	}

	_, err := tailedbeastCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tailedbeast"})
		return
	}

	var updatedTailedBeast struct {
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Images      []string `json:"images"`
		Rank        string   `json:"rank"`
		Abilities   []string `json:"abilities"`
		Personality string   `json:"personality"`
	}

	err = tailedbeastCollection.FindOne(context.TODO(), filter).Decode(&updatedTailedBeast)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": "Successfully updated tailed beast",
		"result":   updatedTailedBeast,
	})
}

// deleteTailedBeast: Menghapus tailedbeast
func DeleteTailedBeast(c *gin.Context) {
	slugParam := c.Param("slug")

	filter := bson.M{"slug": slugParam}
	result, err := tailedbeastCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tailed Beast not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": "Successfully deleted tailed beast",
	})
}
