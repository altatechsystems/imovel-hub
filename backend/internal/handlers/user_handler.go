package handlers

import (
	"net/http"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for administrative users (not brokers)
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
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
