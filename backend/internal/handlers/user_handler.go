package handlers

import (
	"net/http"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/services"
	"github.com/altatech/ecosistema-imob/backend/internal/storage"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for administrative users (not brokers)
type UserHandler struct {
	userService    *services.UserService
	storageService *storage.StorageService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService, storageService *storage.StorageService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		storageService: storageService,
	}
}

// RegisterRoutes registers user routes (tenant-scoped)
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:userId", h.GetUser)
		users.PUT("/:userId", h.UpdateUser)
		users.DELETE("/:userId", h.DeleteUser)
		users.GET("", h.ListUsers)
		users.POST("/:userId/permissions", h.GrantPermission)
		users.DELETE("/:userId/permissions/:permission", h.RevokePermission)
		users.POST("/:userId/photo", h.UploadPhoto)
		users.DELETE("/:userId/photo", h.DeletePhoto)
	}
}

// CreateUser handles POST /api/v1/admin/:tenant_id/users
// Creates a new administrative user (NOT a broker)
func (h *UserHandler) CreateUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.TenantID = tenantID

	if err := h.userService.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /api/v1/admin/:tenant_id/users/:userId
func (h *UserHandler) GetUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")

	user, err := h.userService.GetUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers handles GET /api/v1/admin/:tenant_id/users
// Lists all administrative users for a tenant
func (h *UserHandler) ListUsers(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	// Optional filter for active users only
	activeOnly := c.Query("active") == "true"

	var users []*models.User
	var err error

	if activeOnly {
		users, err = h.userService.ListActiveUsers(c.Request.Context(), tenantID)
	} else {
		users, err = h.userService.ListUsers(c.Request.Context(), tenantID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser handles PUT /api/v1/admin/:tenant_id/users/:userId
func (h *UserHandler) UpdateUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateUser(c.Request.Context(), tenantID, userID, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated user to return
	user, err := h.userService.GetUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User updated but failed to fetch"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /api/v1/admin/:tenant_id/users/:userId
func (h *UserHandler) DeleteUser(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")

	if err := h.userService.DeleteUser(c.Request.Context(), tenantID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GrantPermission handles POST /api/v1/admin/:tenant_id/users/:userId/permissions
// Grants a specific permission to a user
func (h *UserHandler) GrantPermission(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")

	var req struct {
		Permission string `json:"permission" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.GrantPermission(c.Request.Context(), tenantID, userID, req.Permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission granted successfully"})
}

// RevokePermission handles DELETE /api/v1/admin/:tenant_id/users/:userId/permissions/:permission
// Revokes a specific permission from a user
func (h *UserHandler) RevokePermission(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")
	permission := c.Param("permission")

	if err := h.userService.RevokePermission(c.Request.Context(), tenantID, userID, permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked successfully"})
}

// UploadPhoto handles user photo upload
// @Summary Upload user photo
// @Description Upload a profile photo for a user
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param userId path string true "User ID"
// @Param file formData file true "Photo file (JPEG, PNG, or WebP, max 5MB)"
// @Success 200 {object} map[string]interface{} "Photo uploaded successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/admin/{tenant_id}/users/{userId}/photo [post]
func (h *UserHandler) UploadPhoto(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")

	// Check if storage service is available
	if h.storageService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage service is not available"})
		return
	}

	// Get user to verify existence
	user, err := h.userService.GetUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	// Validate file size (max 5MB for profile photos)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 5MB limit"})
		return
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type. Only JPEG, PNG, and WebP are allowed"})
		return
	}

	// Delete old photo if exists
	if user.PhotoURL != "" {
		_ = h.storageService.DeleteBrokerPhoto(c.Request.Context(), tenantID, userID)
	}

	// Upload new photo (reuse broker photo upload logic)
	photoURL, err := h.storageService.UploadBrokerPhoto(c.Request.Context(), tenantID, userID, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user with new photo URL
	if err := h.userService.UpdateUser(c.Request.Context(), tenantID, userID, map[string]interface{}{
		"photo_url": photoURL,
	}); err != nil {
		// Try to delete uploaded photo if update fails
		_ = h.storageService.DeleteBrokerPhoto(c.Request.Context(), tenantID, userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user with photo URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"photo_url": photoURL,
			"message":   "photo uploaded successfully",
		},
	})
}

// DeletePhoto handles user photo deletion
// @Summary Delete user photo
// @Description Delete a user's profile photo
// @Tags users
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{} "Photo deleted successfully"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/admin/{tenant_id}/users/{userId}/photo [delete]
func (h *UserHandler) DeletePhoto(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	userID := c.Param("userId")

	// Check if storage service is available
	if h.storageService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage service is not available"})
		return
	}

	// Get user to verify existence
	user, err := h.userService.GetUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete photo from storage
	if user.PhotoURL != "" {
		if err := h.storageService.DeleteBrokerPhoto(c.Request.Context(), tenantID, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete photo"})
			return
		}
	}

	// Update user to remove photo URL
	if err := h.userService.UpdateUser(c.Request.Context(), tenantID, userID, map[string]interface{}{
		"photo_url": "",
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{"message": "photo deleted successfully"},
	})
}
