package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"cloud.google.com/go/firestore"
)

type AuthHandler struct {
	firebaseAuth *auth.Client
	firestoreDB  *firestore.Client
}

func NewAuthHandler(firebaseAuth *auth.Client, firestoreDB *firestore.Client) *AuthHandler {
	return &AuthHandler{
		firebaseAuth: firebaseAuth,
		firestoreDB:  firestoreDB,
	}
}

// SignupRequest represents the signup payload
type SignupRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	TenantName string `json:"tenant_name" binding:"required"`
}

// SignupResponse represents the signup response
type SignupResponse struct {
	TenantID      string                 `json:"tenant_id"`
	BrokerID      string                 `json:"broker_id"`
	FirebaseToken string                 `json:"firebase_token"`
	User          map[string]interface{} `json:"user"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	FirebaseToken string                 `json:"firebase_token"`
	TenantID      string                 `json:"tenant_id"`
	Broker        map[string]interface{} `json:"broker"`
}

// Signup creates a new tenant and admin broker
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// 1. Check if email already exists
	_, err := h.firebaseAuth.GetUserByEmail(ctx, req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// 2. Create Firebase Auth user
	userParams := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.Name).
		EmailVerified(false)

	userRecord, err := h.firebaseAuth.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating Firebase user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 3. Create Tenant
	tenantID := uuid.New().String()
	slug := generateSlug(req.TenantName)

	tenant := models.Tenant{
		ID:        tenantID,
		Name:      req.TenantName,
		Slug:      slug,
		Email:     req.Email,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.firestoreDB.Collection("tenants").Doc(tenantID).Set(ctx, tenant)
	if err != nil {
		log.Printf("Error creating tenant: %v", err)
		// Rollback: delete Firebase user
		h.firebaseAuth.DeleteUser(ctx, userRecord.UID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
		return
	}

	// 4. Create Broker (admin role)
	brokerID := uuid.New().String()
	broker := models.Broker{
		ID:          brokerID,
		TenantID:    tenantID,
		FirebaseUID: userRecord.UID,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Role:        "admin", // First user is always admin
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = h.firestoreDB.Collection("tenants").Doc(tenantID).Collection("brokers").Doc(brokerID).Set(ctx, broker)
	if err != nil {
		log.Printf("Error creating broker: %v", err)
		// Rollback: delete tenant and Firebase user
		h.firestoreDB.Collection("tenants").Doc(tenantID).Delete(ctx)
		h.firebaseAuth.DeleteUser(ctx, userRecord.UID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create broker"})
		return
	}

	// 5. Set custom claims
	claims := map[string]interface{}{
		"tenant_id": tenantID,
		"role":      "admin",
		"broker_id": brokerID,
	}

	err = h.firebaseAuth.SetCustomUserClaims(ctx, userRecord.UID, claims)
	if err != nil {
		log.Printf("Error setting custom claims: %v", err)
		// Continue anyway, claims can be set later
	}

	// 6. Generate custom token
	token, err := h.firebaseAuth.CustomToken(ctx, userRecord.UID)
	if err != nil {
		log.Printf("Error creating custom token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// 7. Log activity
	go h.logActivity(ctx, tenantID, "tenant_created", map[string]interface{}{
		"tenant_id":   tenantID,
		"tenant_name": req.TenantName,
		"admin_email": req.Email,
	})

	go h.logActivity(ctx, tenantID, "broker_created", map[string]interface{}{
		"broker_id":    brokerID,
		"broker_email": req.Email,
		"role":         "admin",
	})

	// 8. Return response
	c.JSON(http.StatusCreated, SignupResponse{
		TenantID:      tenantID,
		BrokerID:      brokerID,
		FirebaseToken: token,
		User: map[string]interface{}{
			"uid":   userRecord.UID,
			"email": req.Email,
			"name":  req.Name,
		},
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// 1. Get user by email
	userRecord, err := h.firebaseAuth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 2. Find broker by Firebase UID
	brokersQuery := h.firestoreDB.CollectionGroup("brokers").
		Where("firebase_uid", "==", userRecord.UID).
		Limit(1)

	docs, err := brokersQuery.Documents(ctx).GetAll()
	if err != nil || len(docs) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Broker not found"})
		return
	}

	brokerDoc := docs[0]
	var broker models.Broker
	if err := brokerDoc.DataTo(&broker); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load broker data"})
		return
	}

	broker.ID = brokerDoc.Ref.ID

	// 3. Check if broker is active
	if !broker.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Broker account is inactive"})
		return
	}

	// 4. Check if tenant is active
	tenantDoc, err := h.firestoreDB.Collection("tenants").Doc(broker.TenantID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tenant data"})
		return
	}

	var tenant models.Tenant
	if err := tenantDoc.DataTo(&tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tenant data"})
		return
	}

	if !tenant.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tenant account is inactive"})
		return
	}

	// 5. Generate custom token
	token, err := h.firebaseAuth.CustomToken(ctx, userRecord.UID)
	if err != nil {
		log.Printf("Error creating custom token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// 6. Return response
	c.JSON(http.StatusOK, LoginResponse{
		FirebaseToken: token,
		TenantID:      broker.TenantID,
		Broker: map[string]interface{}{
			"id":    broker.ID,
			"name":  broker.Name,
			"email": broker.Email,
			"role":  broker.Role,
		},
	})
}

// RefreshToken refreshes the user's token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()

	// Generate new custom token
	token, err := h.firebaseAuth.CustomToken(ctx, userID.(string))
	if err != nil {
		log.Printf("Error creating custom token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	// Get broker info
	tenantID, _ := c.Get("tenant_id")
	brokerID, _ := c.Get("broker_id")

	c.JSON(http.StatusOK, gin.H{
		"firebase_token": token,
		"tenant_id":      tenantID,
		"broker_id":      brokerID,
	})
}

// Helper: Generate slug from tenant name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only letters, numbers, hyphens)
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	slug = result.String()

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Add timestamp to ensure uniqueness
	timestamp := time.Now().Unix()
	slug = fmt.Sprintf("%s-%d", slug, timestamp)

	return slug
}

// Helper: Log activity
func (h *AuthHandler) logActivity(ctx context.Context, tenantID, eventType string, metadata map[string]interface{}) {
	activityLog := map[string]interface{}{
		"tenant_id":  tenantID,
		"event_type": eventType,
		"metadata":   metadata,
		"timestamp":  time.Now(),
	}

	_, err := h.firestoreDB.Collection("tenants").Doc(tenantID).Collection("activity_logs").NewDoc().Set(ctx, activityLog)
	if err != nil {
		log.Printf("Error logging activity: %v", err)
	}
}
