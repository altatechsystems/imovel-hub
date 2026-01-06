package services

import (
	"context"
	"fmt"
	"time"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/repositories"
	"github.com/altatech/ecosistema-imob/backend/internal/utils"
)

// UserService handles business logic for administrative user management
type UserService struct {
	userRepo        *repositories.UserRepository
	tenantRepo      *repositories.TenantRepository
	activityLogRepo *repositories.ActivityLogRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo *repositories.UserRepository,
	tenantRepo *repositories.TenantRepository,
	activityLogRepo *repositories.ActivityLogRepository,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		tenantRepo:      tenantRepo,
		activityLogRepo: activityLogRepo,
	}
}

// CreateUser creates a new administrative user with validation
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Validate required fields
	if user.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if user.Name == "" {
		return fmt.Errorf("user name is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.FirebaseUID == "" {
		return fmt.Errorf("firebase_uid is required")
	}

	// Validate tenant exists
	if _, err := s.tenantRepo.Get(ctx, user.TenantID); err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	// Validate email
	if err := utils.ValidateEmail(user.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	user.Email = utils.NormalizeEmail(user.Email)

	// Check if email is already registered in this tenant
	existing, err := s.userRepo.GetByEmail(ctx, user.TenantID, user.Email)
	if err != nil && err != repositories.ErrNotFound {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("email '%s' is already registered in this tenant", user.Email)
	}

	// Check if Firebase UID is already registered in this tenant
	existing, err = s.userRepo.GetByFirebaseUID(ctx, user.TenantID, user.FirebaseUID)
	if err != nil && err != repositories.ErrNotFound {
		return fmt.Errorf("failed to check firebase_uid uniqueness: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("firebase_uid '%s' is already registered in this tenant", user.FirebaseUID)
	}

	// Validate phone if provided
	if user.Phone != "" {
		// Remove country code if present (+55)
		cleanPhone := user.Phone
		if len(cleanPhone) > 11 {
			// Try to strip +55 or 55
			if cleanPhone[:3] == "+55" {
				cleanPhone = cleanPhone[3:]
			} else if cleanPhone[:2] == "55" {
				cleanPhone = cleanPhone[2:]
			}
		}

		if err := utils.ValidatePhoneBR(cleanPhone); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
		user.Phone = utils.NormalizePhoneBR(cleanPhone)
	}

	// Validate document if provided
	if user.Document != "" {
		if user.DocumentType == "" {
			return fmt.Errorf("document_type is required when document is provided")
		}

		if user.DocumentType == "cpf" {
			if err := utils.ValidateCPF(user.Document); err != nil {
				return fmt.Errorf("invalid CPF: %w", err)
			}
			user.Document = utils.NormalizeCPF(user.Document)
		} else if user.DocumentType == "cnpj" {
			if err := utils.ValidateCNPJ(user.Document); err != nil {
				return fmt.Errorf("invalid CNPJ: %w", err)
			}
			user.Document = utils.NormalizeCNPJ(user.Document)
		} else {
			return fmt.Errorf("document_type must be 'cpf' or 'cnpj'")
		}
	}

	// Validate role
	if user.Role == "" {
		user.Role = "admin" // Default role
	}
	if !models.IsValidUserRole(user.Role) {
		return fmt.Errorf("invalid role: must be 'admin' or 'manager'")
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Set active by default
	if !user.IsActive {
		user.IsActive = true
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Log activity
	go s.logActivity(ctx, user.TenantID, "user_created", models.ActorTypeSystem, "system", map[string]interface{}{
		"user_id":    user.ID,
		"user_email": user.Email,
		"role":       user.Role,
	})

	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, tenantID, userID string, updates map[string]interface{}) error {
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validate user exists
	existing, err := s.userRepo.Get(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Validate email if being updated
	if email, ok := updates["email"].(string); ok && email != "" {
		if err := utils.ValidateEmail(email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
		normalizedEmail := utils.NormalizeEmail(email)

		// Check if email is already registered by another user in this tenant
		emailUser, err := s.userRepo.GetByEmail(ctx, tenantID, normalizedEmail)
		if err != nil && err != repositories.ErrNotFound {
			return fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if emailUser != nil && emailUser.ID != userID {
			return fmt.Errorf("email '%s' is already registered by another user", normalizedEmail)
		}

		updates["email"] = normalizedEmail
	}

	// Validate phone if being updated
	if phone, ok := updates["phone"].(string); ok && phone != "" {
		cleanPhone := phone
		if len(cleanPhone) > 11 {
			if cleanPhone[:3] == "+55" {
				cleanPhone = cleanPhone[3:]
			} else if cleanPhone[:2] == "55" {
				cleanPhone = cleanPhone[2:]
			}
		}

		if err := utils.ValidatePhoneBR(cleanPhone); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
		updates["phone"] = utils.NormalizePhoneBR(cleanPhone)
	}

	// Validate document if being updated
	if document, ok := updates["document"].(string); ok && document != "" {
		documentType, hasType := updates["document_type"].(string)
		if !hasType || documentType == "" {
			documentType = existing.DocumentType
		}

		if documentType == "" {
			return fmt.Errorf("document_type is required when updating document")
		}

		if documentType == "cpf" {
			if err := utils.ValidateCPF(document); err != nil {
				return fmt.Errorf("invalid CPF: %w", err)
			}
			updates["document"] = utils.NormalizeCPF(document)
			updates["document_type"] = "cpf"
		} else if documentType == "cnpj" {
			if err := utils.ValidateCNPJ(document); err != nil {
				return fmt.Errorf("invalid CNPJ: %w", err)
			}
			updates["document"] = utils.NormalizeCNPJ(document)
			updates["document_type"] = "cnpj"
		} else {
			return fmt.Errorf("document_type must be 'cpf' or 'cnpj'")
		}
	}

	// Validate role if being updated
	if role, ok := updates["role"].(string); ok && role != "" {
		if !models.IsValidUserRole(role) {
			return fmt.Errorf("invalid role: must be 'admin' or 'manager'")
		}
	}

	// Update timestamps
	updates["updated_at"] = time.Now()

	// Update user
	if err := s.userRepo.Update(ctx, tenantID, userID, updates); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Log activity
	go s.logActivity(ctx, tenantID, "user_updated", models.ActorTypeSystem, "system", map[string]interface{}{
		"user_id": userID,
		"updates": updates,
	})

	return nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, tenantID, userID string) (*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	user, err := s.userRepo.Get(ctx, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ListUsers retrieves all users for a tenant
func (s *UserService) ListUsers(ctx context.Context, tenantID string) ([]*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	users, err := s.userRepo.List(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// ListActiveUsers retrieves all active users for a tenant
func (s *UserService) ListActiveUsers(ctx context.Context, tenantID string) ([]*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	users, err := s.userRepo.ListActive(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active users: %w", err)
	}

	return users, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, tenantID, userID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Check if user exists
	user, err := s.userRepo.Get(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, tenantID, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Log activity
	go s.logActivity(ctx, tenantID, "user_deleted", models.ActorTypeSystem, "system", map[string]interface{}{
		"user_id":    userID,
		"user_email": user.Email,
	})

	return nil
}

// GrantPermission adds a permission to a user
func (s *UserService) GrantPermission(ctx context.Context, tenantID, userID, permission string) error {
	user, err := s.userRepo.Get(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Add permission
	user.AddPermission(permission)

	// Update user
	updates := map[string]interface{}{
		"permissions": user.Permissions,
		"updated_at":  time.Now(),
	}

	if err := s.userRepo.Update(ctx, tenantID, userID, updates); err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	// Log activity
	go s.logActivity(ctx, tenantID, "permission_granted", models.ActorTypeSystem, "system", map[string]interface{}{
		"user_id":    userID,
		"permission": permission,
	})

	return nil
}

// RevokePermission removes a permission from a user
func (s *UserService) RevokePermission(ctx context.Context, tenantID, userID, permission string) error {
	user, err := s.userRepo.Get(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Remove permission
	user.RemovePermission(permission)

	// Update user
	updates := map[string]interface{}{
		"permissions": user.Permissions,
		"updated_at":  time.Now(),
	}

	if err := s.userRepo.Update(ctx, tenantID, userID, updates); err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	// Log activity
	go s.logActivity(ctx, tenantID, "permission_revoked", models.ActorTypeSystem, "system", map[string]interface{}{
		"user_id":    userID,
		"permission": permission,
	})

	return nil
}

// logActivity logs an activity (helper method)
func (s *UserService) logActivity(ctx context.Context, tenantID, eventType string, actorType models.ActorType, actorID string, metadata map[string]interface{}) error {
	log := &models.ActivityLog{
		TenantID:  tenantID,
		EventType: eventType,
		ActorType: actorType,
		ActorID:   actorID,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	return s.activityLogRepo.Create(ctx, log)
}
