package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/repositories"
	"github.com/altatech/ecosistema-imob/backend/internal/utils"
)

// TenantService handles business logic for tenant management
type TenantService struct {
	tenantRepo      *repositories.TenantRepository
	activityLogRepo *repositories.ActivityLogRepository
}

// NewTenantService creates a new tenant service
func NewTenantService(
	tenantRepo *repositories.TenantRepository,
	activityLogRepo *repositories.ActivityLogRepository,
) *TenantService {
	return &TenantService{
		tenantRepo:      tenantRepo,
		activityLogRepo: activityLogRepo,
	}
}

// CreateTenant creates a new tenant with validation
func (s *TenantService) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	// Validate required fields
	if tenant.Name == "" {
		return fmt.Errorf("tenant name is required")
	}

	// Generate slug from name if not provided
	if tenant.Slug == "" {
		tenant.Slug = s.GenerateSlug(tenant.Name)
	} else {
		// Normalize provided slug
		tenant.Slug = s.NormalizeSlug(tenant.Slug)
	}

	// Validate slug uniqueness
	if err := s.ValidateSlug(ctx, tenant.Slug, ""); err != nil {
		return err
	}

	// Validate Document (CPF or CNPJ) if provided
	if tenant.Document != "" {
		// Try to detect document type if not specified
		cleanDoc := regexp.MustCompile(`[^\d]`).ReplaceAllString(tenant.Document, "")

		if tenant.DocumentType == "" {
			// Auto-detect based on length
			if len(cleanDoc) == 11 {
				tenant.DocumentType = "cpf"
			} else if len(cleanDoc) == 14 {
				tenant.DocumentType = "cnpj"
			} else {
				return fmt.Errorf("document must be CPF (11 digits) or CNPJ (14 digits)")
			}
		}

		// Validate based on type
		if tenant.DocumentType == "cpf" {
			if err := utils.ValidateCPF(tenant.Document); err != nil {
				return fmt.Errorf("invalid CPF: %w", err)
			}
			tenant.Document = utils.NormalizeCPF(tenant.Document)
		} else if tenant.DocumentType == "cnpj" {
			if err := utils.ValidateCNPJ(tenant.Document); err != nil {
				return fmt.Errorf("invalid CNPJ: %w", err)
			}
			tenant.Document = utils.NormalizeCNPJ(tenant.Document)
		} else {
			return fmt.Errorf("document_type must be 'cpf' or 'cnpj'")
		}
	}

	// Validate business_type if provided
	if tenant.BusinessType != "" {
		validBusinessTypes := map[string]bool{
			"imobiliaria":       true,
			"incorporadora":     true,
			"loteadora":         true,
			"construtora":       true,
			"corretor_autonomo": true,
		}
		if !validBusinessTypes[tenant.BusinessType] {
			return fmt.Errorf("business_type must be one of: imobiliaria, incorporadora, loteadora, construtora, corretor_autonomo")
		}
	}

	// Validate CRECI if provided
	if tenant.CRECI != "" {
		if err := utils.ValidateCRECI(tenant.CRECI); err != nil {
			return fmt.Errorf("invalid CRECI: %w", err)
		}
		tenant.CRECI = utils.NormalizeCRECI(tenant.CRECI)
	}

	// Validate email if provided
	if tenant.Email != "" {
		if err := utils.ValidateEmail(tenant.Email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
		tenant.Email = utils.NormalizeEmail(tenant.Email)
	}

	// Validate phone if provided
	if tenant.Phone != "" {
		// Remove country code if present (+55)
		cleanPhone := strings.TrimPrefix(tenant.Phone, "+55")
		cleanPhone = strings.TrimPrefix(cleanPhone, "55")

		if err := utils.ValidatePhoneBR(cleanPhone); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
		tenant.Phone = utils.NormalizePhoneBR(cleanPhone)
	}

	// Set defaults
	if tenant.Country == "" {
		tenant.Country = "BR"
	}

	// Set active by default
	tenant.IsActive = true

	// Create tenant in repository
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	// Log activity
	_ = s.logActivity(ctx, tenant.ID, "tenant_created", models.ActorTypeSystem, "", map[string]interface{}{
		"tenant_id": tenant.ID,
		"name":      tenant.Name,
		"slug":      tenant.Slug,
	})

	return nil
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(ctx context.Context, id string) (*models.Tenant, error) {
	if id == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}

	tenant, err := s.tenantRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return tenant, nil
}

// GetTenantBySlug retrieves a tenant by slug
func (s *TenantService) GetTenantBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	if slug == "" {
		return nil, fmt.Errorf("tenant slug is required")
	}

	slug = s.NormalizeSlug(slug)

	tenant, err := s.tenantRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by slug: %w", err)
	}

	return tenant, nil
}

// UpdateTenant updates a tenant with validation
func (s *TenantService) UpdateTenant(ctx context.Context, id string, updates map[string]interface{}) error {
	if id == "" {
		return fmt.Errorf("tenant ID is required")
	}

	// Validate tenant exists
	existing, err := s.tenantRepo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	// Validate slug if being updated
	if slug, ok := updates["slug"].(string); ok {
		normalized := s.NormalizeSlug(slug)
		if err := s.ValidateSlug(ctx, normalized, id); err != nil {
			return err
		}
		updates["slug"] = normalized
	}

	// Validate Document (CPF or CNPJ) if being updated
	if document, ok := updates["document"].(string); ok && document != "" {
		// Get document type from updates or existing tenant
		documentType, hasType := updates["document_type"].(string)
		if !hasType || documentType == "" {
			documentType = existing.DocumentType
		}

		// Try to detect document type if not specified
		cleanDoc := regexp.MustCompile(`[^\d]`).ReplaceAllString(document, "")

		if documentType == "" {
			// Auto-detect based on length
			if len(cleanDoc) == 11 {
				documentType = "cpf"
			} else if len(cleanDoc) == 14 {
				documentType = "cnpj"
			} else {
				return fmt.Errorf("document must be CPF (11 digits) or CNPJ (14 digits)")
			}
		}

		// Validate based on type
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

	// Validate business_type if being updated
	if businessType, ok := updates["business_type"].(string); ok && businessType != "" {
		validBusinessTypes := map[string]bool{
			"imobiliaria":       true,
			"incorporadora":     true,
			"loteadora":         true,
			"construtora":       true,
			"corretor_autonomo": true,
		}
		if !validBusinessTypes[businessType] {
			return fmt.Errorf("business_type must be one of: imobiliaria, incorporadora, loteadora, construtora, corretor_autonomo")
		}
	}

	// Validate CRECI if being updated
	if creci, ok := updates["creci"].(string); ok && creci != "" {
		if err := utils.ValidateCRECI(creci); err != nil {
			return fmt.Errorf("invalid CRECI: %w", err)
		}
		updates["creci"] = utils.NormalizeCRECI(creci)
	}

	// Validate email if being updated
	if email, ok := updates["email"].(string); ok && email != "" {
		if err := utils.ValidateEmail(email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
		updates["email"] = utils.NormalizeEmail(email)
	}

	// Validate phone if being updated
	if phone, ok := updates["phone"].(string); ok && phone != "" {
		// Remove country code if present (+55)
		cleanPhone := strings.TrimPrefix(phone, "+55")
		cleanPhone = strings.TrimPrefix(cleanPhone, "55")

		if err := utils.ValidatePhoneBR(cleanPhone); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
		updates["phone"] = utils.NormalizePhoneBR(cleanPhone)
	}

	// Update tenant in repository
	if err := s.tenantRepo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	// Log activity
	_ = s.logActivity(ctx, id, "tenant_updated", models.ActorTypeSystem, "", map[string]interface{}{
		"tenant_id": id,
		"updates":   updates,
	})

	_ = existing // silence unused warning

	return nil
}

// DeleteTenant deletes a tenant
func (s *TenantService) DeleteTenant(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("tenant ID is required")
	}

	// Validate tenant exists
	if _, err := s.tenantRepo.Get(ctx, id); err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	// Delete tenant from repository
	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	// Log activity
	_ = s.logActivity(ctx, id, "tenant_deleted", models.ActorTypeSystem, "", map[string]interface{}{
		"tenant_id": id,
	})

	return nil
}

// ListTenants lists all tenants with pagination
func (s *TenantService) ListTenants(ctx context.Context, opts repositories.PaginationOptions) ([]*models.Tenant, error) {
	tenants, err := s.tenantRepo.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return tenants, nil
}

// ListActiveTenants lists all active tenants
func (s *TenantService) ListActiveTenants(ctx context.Context, opts repositories.PaginationOptions) ([]*models.Tenant, error) {
	tenants, err := s.tenantRepo.ListActive(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list active tenants: %w", err)
	}

	return tenants, nil
}

// ActivateTenant activates a tenant
func (s *TenantService) ActivateTenant(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("tenant ID is required")
	}

	updates := map[string]interface{}{
		"is_active": true,
	}

	if err := s.tenantRepo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to activate tenant: %w", err)
	}

	// Log activity
	_ = s.logActivity(ctx, id, "tenant_activated", models.ActorTypeSystem, "", map[string]interface{}{
		"tenant_id": id,
	})

	return nil
}

// DeactivateTenant deactivates a tenant
func (s *TenantService) DeactivateTenant(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("tenant ID is required")
	}

	updates := map[string]interface{}{
		"is_active": false,
	}

	if err := s.tenantRepo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to deactivate tenant: %w", err)
	}

	// Log activity
	_ = s.logActivity(ctx, id, "tenant_deactivated", models.ActorTypeSystem, "", map[string]interface{}{
		"tenant_id": id,
	})

	return nil
}

// ValidateSlug validates slug uniqueness
// Pass the current tenant ID to exclude it from uniqueness check (for updates)
func (s *TenantService) ValidateSlug(ctx context.Context, slug string, currentTenantID string) error {
	if slug == "" {
		return fmt.Errorf("slug is required")
	}

	// Check if slug is already taken
	existing, err := s.tenantRepo.GetBySlug(ctx, slug)
	if err != nil {
		// If not found, slug is available
		if err == repositories.ErrNotFound {
			return nil
		}
		return fmt.Errorf("failed to validate slug: %w", err)
	}

	// If found and it's a different tenant, slug is taken
	if existing.ID != currentTenantID {
		return fmt.Errorf("slug '%s' is already taken", slug)
	}

	return nil
}

// GenerateSlug generates a URL-friendly slug from a name
func (s *TenantService) GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Remove accents and special characters
	slug = s.removeAccents(slug)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 50 {
		slug = slug[:50]
	}

	return slug
}

// NormalizeSlug normalizes a slug
func (s *TenantService) NormalizeSlug(slug string) string {
	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Remove accents and special characters
	slug = s.removeAccents(slug)

	// Replace invalid characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Replace multiple consecutive hyphens with single hyphen
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// removeAccents removes accents from a string using Unicode normalization
func (s *TenantService) removeAccents(str string) string {
	// Use NFD (Normalization Form Decomposed) to separate base characters from diacritics
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		// Remove combining diacritical marks (accents)
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)

	result, _, err := transform.String(t, str)
	if err != nil {
		// Fallback to original string if transformation fails
		return str
	}

	return result
}

// logActivity logs an activity (helper method)
func (s *TenantService) logActivity(ctx context.Context, tenantID, eventType string, actorType models.ActorType, actorID string, metadata map[string]interface{}) error {
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
