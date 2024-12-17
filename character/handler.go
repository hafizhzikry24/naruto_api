package character

import (
	"net/http"
	"strconv"

	"my-gin-app/models"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// CreateUser handler untuk membuat karakter baru
func (h *Handler) CreateUser(c *gin.Context) {
	var character models.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateCharacter(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": character})
}

// ReadUser handler untuk membaca karakter berdasarkan slug
func (h *Handler) ReadUser(c *gin.Context) {
	slugParam := c.Param("slug")
	character, err := h.Service.GetCharacterBySlug(slugParam)
	if err != nil {
		if err.Error() == "character not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success retrieved data",
		"result":  character,
	})
}

// UpdateUser handler untuk memperbarui karakter berdasarkan slug
func (h *Handler) UpdateUser(c *gin.Context) {
	slugParam := c.Param("slug")
	var updatedData models.Character
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateCharacter(slugParam, &updatedData); err != nil {
		if err.Error() == "character not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	updatedCharacter, err := h.Service.GetCharacterBySlug(slug.Make(updatedData.Name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated character"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Character updated",
		"result":  updatedCharacter,
	})
}

// DeleteUser handler untuk menghapus karakter berdasarkan slug
func (h *Handler) DeleteUser(c *gin.Context) {
	slugParam := c.Param("slug")
	if err := h.Service.DeleteCharacter(slugParam); err != nil {
		if err.Error() == "character not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Character Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Character deleted"})
}

// IndexUser handler untuk mengambil daftar karakter dengan pagination
func (h *Handler) IndexUser(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var page, limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
			return
		}
	}

	characters, count, err := h.Service.ListCharacters(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if page > 0 && limit > 0 {
		totalPages := (count + int64(limit) - 1) / int64(limit)
		c.JSON(http.StatusOK, gin.H{
			"message":    "Success retrieved data",
			"result":     characters,
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
			"totalItems": count,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Success retrieved all data",
			"result":  characters,
		})
	}
}

// SearchCharacter handler untuk mencari karakter berdasarkan nama
func (h *Handler) SearchCharacter(c *gin.Context) {
	nameQuery := c.Query("name")
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
		return
	}

	characters, err := h.Service.SearchCharacters(nameQuery)
	if err != nil {
		if err.Error() == "no characters found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No characters found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Found characters",
		"result":  characters,
	})
}
