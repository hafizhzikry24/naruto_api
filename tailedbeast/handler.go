package tailedbeast

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

// CreateTailedBeast handler untuk menambahkan tailedbeast baru
func (h *Handler) CreateTailedBeast(c *gin.Context) {
	var beast models.TailedBeast
	if err := c.ShouldBindJSON(&beast); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateBeast(&beast); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": beast})
}

// ReadTailedBeast handler untuk membaca tailedbeast berdasarkan slug
func (h *Handler) ReadTailedBeast(c *gin.Context) {
	slugParam := c.Param("slug")
	beast, err := h.Service.GetBeastBySlug(slugParam)
	if err != nil {
		if err.Error() == "tailed beast not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tailed Beast Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success retrieved data",
		"result":  beast,
	})
}

// UpdateTailedBeast handler untuk memperbarui tailedbeast berdasarkan slug
func (h *Handler) UpdateTailedBeast(c *gin.Context) {
	slugParam := c.Param("slug")
	var updatedData models.TailedBeast
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateBeast(slugParam, &updatedData); err != nil {
		if err.Error() == "tailed beast not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tailed Beast Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	updatedBeast, err := h.Service.GetBeastBySlug(slug.Make(updatedData.Name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated tailed beast"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Tailed Beast updated",
		"result":  updatedBeast,
	})
}

// DeleteTailedBeast handler untuk menghapus tailedbeast berdasarkan slug
func (h *Handler) DeleteTailedBeast(c *gin.Context) {
	slugParam := c.Param("slug")
	if err := h.Service.DeleteBeast(slugParam); err != nil {
		if err.Error() == "tailed beast not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tailed Beast Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Tailed Beast deleted"})
}

// IndexTailedBeast handler untuk mengambil daftar tailedbeast dengan pagination
func (h *Handler) IndexTailedBeast(c *gin.Context) {
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

	beasts, count, err := h.Service.ListBeasts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if page > 0 && limit > 0 {
		totalPages := (count + int64(limit) - 1) / int64(limit)
		c.JSON(http.StatusOK, gin.H{
			"message":    "Success retrieved data",
			"result":     beasts,
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
			"totalItems": count,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Success retrieved all data",
			"result":  beasts,
		})
	}
}

// SearchTailedBeast handler untuk mencari tailedbeast berdasarkan nama
func (h *Handler) SearchTailedBeast(c *gin.Context) {
	nameQuery := c.Query("name")
	if nameQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
		return
	}

	beasts, err := h.Service.SearchBeasts(nameQuery)
	if err != nil {
		if err.Error() == "no tailed beasts found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No tailed beasts found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Found tailed beasts",
		"result":  beasts,
	})
}
