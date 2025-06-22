package handlers

import (
	"fmt"
	"my-property/go-service/database"
	"my-property/go-service/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadPropertyImages(c *gin.Context) {
	propertyID := c.Param("id")

	var property models.Building
	if err := database.DB.First(&property, "id = ?", propertyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("get form err: %s", err.Error())})
		return
	}
	files := form.File["images"]

	// Create storage directory if it doesn't exist
	storagePath := filepath.Join("storage", "properties", propertyID)
	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}

	var newImages []models.PropertyImage
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		// Use a public-facing URL path
		urlPath := filepath.Join("/storage", "properties", propertyID, filename)
		// And a corresponding filesystem path to save the file
		fsPath := filepath.Join("storage", "properties", propertyID, filename)

		if err := c.SaveUploadedFile(file, fsPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("upload file err: %s", err.Error())})
			return
		}

		newImage := models.PropertyImage{
			PropertyID: property.ID,
			URL:        urlPath,
		}
		newImages = append(newImages, newImage)
	}

	if err := database.DB.Create(&newImages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image records to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"images":  newImages,
	})
}
