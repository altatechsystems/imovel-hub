package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/altatech/ecosistema-imob/backend/internal/services"
	"github.com/altatech/ecosistema-imob/backend/internal/storage"
	"github.com/gin-gonic/gin"
)

// multipartFileWrapper wraps bytes.Reader to implement multipart.File
type multipartFileWrapper struct {
	*bytes.Reader
}

func (m *multipartFileWrapper) Close() error {
	return nil
}

// bytesToMultipartFile converts bytes to multipart.File
func bytesToMultipartFile(data []byte) multipart.File {
	return &multipartFileWrapper{bytes.NewReader(data)}
}

// StorageHandler handles storage-related HTTP requests
type StorageHandler struct {
	storageService *storage.StorageService
	photoProcessor *services.PhotoProcessor
}

// NewStorageHandler creates a new storage handler
func NewStorageHandler(storageService *storage.StorageService, photoProcessor *services.PhotoProcessor) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
		photoProcessor: photoProcessor,
	}
}

// RegisterRoutes registers storage routes (tenant-scoped)
func (h *StorageHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Using /property-images as a separate route to avoid conflict with /properties/:id
	images := router.Group("/property-images/:property_id")
	{
		images.POST("", h.UploadImage)
		images.GET("", h.ListImages)
		images.GET("/:image_id", h.GetImageURL)
		images.DELETE("/:image_id", h.DeleteImage)
	}
}

// UploadImage handles image upload with automatic photo processing
// @Summary Upload property image
// @Description Upload an image for a property (automatically processes into 3 sizes: thumb, medium, large)
// @Tags storage
// @Accept multipart/form-data
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param property_id path string true "Property ID"
// @Param file formData file true "Image file (JPEG, PNG, or WebP, max 10MB)"
// @Param order formData int false "Display order (0-based, default: 0)"
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

	// Get optional order parameter
	order := 0
	if orderStr := c.PostForm("order"); orderStr != "" {
		var orderVal int
		if _, err := fmt.Sscanf(orderStr, "%d", &orderVal); err == nil {
			order = orderVal
		}
	}

	// Check if photo processor is available
	if h.photoProcessor == nil {
		// Fallback to simple upload without processing
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
		return
	}

	// Read file into memory for processing
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to read uploaded file",
		})
		return
	}

	// Create a temporary HTTP server to serve the file for PhotoProcessor
	// PhotoProcessor expects a URL, so we need to upload to GCS first, then process
	// Alternative: Modify PhotoProcessor to accept byte data directly
	// For now, upload original file first, then process it

	// Upload original file temporarily
	originalMetadata, err := h.storageService.UploadPropertyImage(
		c.Request.Context(),
		tenantID,
		propertyID,
		bytesToMultipartFile(fileBytes),
		header,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Process photo using the uploaded URL
	photo, err := h.photoProcessor.ProcessPhoto(
		c.Request.Context(),
		tenantID,
		propertyID,
		originalMetadata.URL,
		order,
	)
	if err != nil {
		// Photo processing failed, but we have the original upload
		// Return the original upload info
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    originalMetadata,
			"warning": "Photo uploaded but processing failed: " + err.Error(),
		})
		return
	}

	// Return processed photo info
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"id":         photo.ID,
			"thumb_url":  photo.ThumbURL,
			"medium_url": photo.MediumURL,
			"large_url":  photo.LargeURL,
			"url":        photo.URL,
			"order":      photo.Order,
			"is_cover":   photo.IsCover,
		},
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
