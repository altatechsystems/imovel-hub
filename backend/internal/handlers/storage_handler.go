package handlers

import (
	"net/http"

	"github.com/altatech/ecosistema-imob/backend/internal/storage"
	"github.com/gin-gonic/gin"
)

// StorageHandler handles storage-related HTTP requests
type StorageHandler struct {
	storageService *storage.StorageService
}

// NewStorageHandler creates a new storage handler
func NewStorageHandler(storageService *storage.StorageService) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
	}
}

// RegisterRoutes registers storage routes (tenant-scoped)
func (h *StorageHandler) RegisterRoutes(router *gin.RouterGroup) {
	images := router.Group("/:tenant_id/properties/:property_id/images")
	{
		images.POST("", h.UploadImage)
		images.GET("", h.ListImages)
		images.GET("/:image_id", h.GetImageURL)
		images.DELETE("/:image_id", h.DeleteImage)
	}
}

// UploadImage handles image upload
// @Summary Upload property image
// @Description Upload an image for a property
// @Tags storage
// @Accept multipart/form-data
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param property_id path string true "Property ID"
// @Param file formData file true "Image file (JPEG, PNG, or WebP, max 10MB)"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 413 {object} map[string]interface{}
// @Failure 415 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/{tenant_id}/properties/{property_id}/images [post]
func (h *StorageHandler) UploadImage(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	propertyID := c.Param("property_id")

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "file is required",
		})
		return
	}
	defer file.Close()

	// Validate file size before processing
	if err := h.storageService.ValidateFileSize(header.Size); err != nil {
		if err == storage.ErrFileTooLarge {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if err := h.storageService.ValidateContentType(contentType); err != nil {
		if err == storage.ErrInvalidFileType {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Upload image
	metadata, err := h.storageService.UploadPropertyImage(
		c.Request.Context(),
		tenantID,
		propertyID,
		file,
		header,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    metadata,
	})
}

// ListImages lists all images for a property
// @Summary List property images
// @Description List all images for a property
// @Tags storage
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param property_id path string true "Property ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/{tenant_id}/properties/{property_id}/images [get]
func (h *StorageHandler) ListImages(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	propertyID := c.Param("property_id")

	images, err := h.storageService.ListPropertyImages(
		c.Request.Context(),
		tenantID,
		propertyID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    images,
		"count":   len(images),
	})
}

// GetImageURL retrieves a signed URL for an image
// @Summary Get image URL
// @Description Get a signed URL for accessing an image
// @Tags storage
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param property_id path string true "Property ID"
// @Param image_id path string true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/{tenant_id}/properties/{property_id}/images/{image_id} [get]
func (h *StorageHandler) GetImageURL(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	propertyID := c.Param("property_id")
	imageID := c.Param("image_id")

	url, err := h.storageService.GetPropertyImageURL(
		c.Request.Context(),
		tenantID,
		propertyID,
		imageID,
	)
	if err != nil {
		if err == storage.ErrImageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "image not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url":       url,
			"image_id":  imageID,
			"expires":   "1 hour",
		},
	})
}

// DeleteImage deletes an image
// @Summary Delete image
// @Description Delete an image from storage
// @Tags storage
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param property_id path string true "Property ID"
// @Param image_id path string true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/{tenant_id}/properties/{property_id}/images/{image_id} [delete]
func (h *StorageHandler) DeleteImage(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	propertyID := c.Param("property_id")
	imageID := c.Param("image_id")

	err := h.storageService.DeletePropertyImage(
		c.Request.Context(),
		tenantID,
		propertyID,
		imageID,
	)
	if err != nil {
		if err == storage.ErrImageNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "image not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{"message": "image deleted successfully"},
	})
}
