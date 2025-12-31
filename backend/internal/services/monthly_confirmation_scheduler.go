package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/repositories"
)

// MonthlyConfirmationScheduler handles automatic monthly confirmation reminders
type MonthlyConfirmationScheduler struct {
	scheduledConfirmationRepo *repositories.ScheduledConfirmationRepository
	propertyRepo              *repositories.PropertyRepository
	ownerRepo                 *repositories.OwnerRepository
	ownerConfirmationService  *OwnerConfirmationService
}

// NewMonthlyConfirmationScheduler creates a new monthly confirmation scheduler
func NewMonthlyConfirmationScheduler(
	scheduledConfirmationRepo *repositories.ScheduledConfirmationRepository,
	propertyRepo *repositories.PropertyRepository,
	ownerRepo *repositories.OwnerRepository,
	ownerConfirmationService *OwnerConfirmationService,
) *MonthlyConfirmationScheduler {
	return &MonthlyConfirmationScheduler{
		scheduledConfirmationRepo: scheduledConfirmationRepo,
		propertyRepo:              propertyRepo,
		ownerRepo:                 ownerRepo,
		ownerConfirmationService:  ownerConfirmationService,
	}
}

// ScheduleMonthlyConfirmationsRequest represents the request to schedule monthly confirmations
type ScheduleMonthlyConfirmationsRequest struct {
	TenantID     string    `json:"tenant_id"`
	ScheduledFor time.Time `json:"scheduled_for"` // When to send (default: 1st of next month)
	DryRun       bool      `json:"dry_run"`       // If true, only returns count without creating
}

// ScheduleMonthlyConfirmationsResponse represents the response
type ScheduleMonthlyConfirmationsResponse struct {
	TotalProperties     int      `json:"total_properties"`
	ScheduledCount      int      `json:"scheduled_count"`
	SkippedCount        int      `json:"skipped_count"`
	SkippedReasons      []string `json:"skipped_reasons,omitempty"`
	ScheduledForDate    string   `json:"scheduled_for_date"`
	ScheduledConfirmIDs []string `json:"scheduled_confirm_ids,omitempty"`
}

// ScheduleMonthlyConfirmations creates scheduled confirmations for all active properties
// This should be run monthly (e.g., on the 1st of each month)
func (s *MonthlyConfirmationScheduler) ScheduleMonthlyConfirmations(ctx context.Context, req ScheduleMonthlyConfirmationsRequest) (*ScheduleMonthlyConfirmationsResponse, error) {
	if req.TenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	// Default to 1st of next month at 9 AM if not specified
	if req.ScheduledFor.IsZero() {
		now := time.Now()
		req.ScheduledFor = time.Date(now.Year(), now.Month()+1, 1, 9, 0, 0, 0, now.Location())
	}

	log.Printf("🗓️  Scheduling monthly confirmations for tenant %s on %s", req.TenantID, req.ScheduledFor.Format("2006-01-02"))

	// Get all properties for the tenant
	properties, err := s.propertyRepo.List(ctx, req.TenantID, &repositories.PropertyFilters{}, repositories.PaginationOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list properties: %w", err)
	}

	response := &ScheduleMonthlyConfirmationsResponse{
		TotalProperties:  len(properties),
		ScheduledForDate: req.ScheduledFor.Format("2006-01-02 15:04:05"),
		SkippedReasons:   []string{},
	}

	// Check if confirmations already exist for this month
	existingConfirmations, err := s.scheduledConfirmationRepo.GetByPropertyAndMonth(
		ctx,
		req.TenantID,
		"", // Will check per property
		req.ScheduledFor.Year(),
		req.ScheduledFor.Month(),
	)
	if err != nil {
		log.Printf("⚠️  Warning: failed to check existing confirmations: %v", err)
	}

	// Create a map for quick lookup
	existingByProperty := make(map[string]bool)
	for _, ec := range existingConfirmations {
		existingByProperty[ec.PropertyID] = true
	}

	for _, property := range properties {
		// Skip if property has no owner
		if property.OwnerID == "" {
			response.SkippedCount++
			response.SkippedReasons = append(response.SkippedReasons, fmt.Sprintf("Property %s: no owner", property.Reference))
			continue
		}

		// Skip if already scheduled for this month
		if existingByProperty[property.ID] {
			response.SkippedCount++
			response.SkippedReasons = append(response.SkippedReasons, fmt.Sprintf("Property %s: already scheduled", property.Reference))
			continue
		}

		// Skip if property status is not available or pending_confirmation
		if property.Status != models.PropertyStatusAvailable && property.Status != models.PropertyStatusPendingConfirmation {
			response.SkippedCount++
			response.SkippedReasons = append(response.SkippedReasons, fmt.Sprintf("Property %s: status is %s", property.Reference, property.Status))
			continue
		}

		if req.DryRun {
			response.ScheduledCount++
			continue
		}

		// Generate confirmation token and link
		confirmationURL, tokenID, expiresAt, err := s.ownerConfirmationService.GenerateOwnerConfirmationLink(
			ctx,
			req.TenantID,
			property.ID,
			property.CaptadorID, // actorID is the broker who captured the property
			&property.OwnerID,   // ownerID as pointer
			"whatsapp",
		)
		if err != nil {
			log.Printf("❌ Failed to generate confirmation link for property %s: %v", property.Reference, err)
			response.SkippedCount++
			response.SkippedReasons = append(response.SkippedReasons, fmt.Sprintf("Property %s: failed to generate link", property.Reference))
			continue
		}

		// Create scheduled confirmation record
		scheduledConfirmation := &models.ScheduledConfirmation{
			TenantID:        req.TenantID,
			PropertyID:      property.ID,
			OwnerID:         property.OwnerID,
			BrokerID:        property.CaptadorID,
			TokenID:         tokenID,
			ConfirmationURL: confirmationURL,
			ScheduledFor:    req.ScheduledFor,
			Status:          models.ScheduledConfirmationStatusPending,
			DeliveryMethod:  "manual", // Will be updated when WhatsApp API is integrated
		}

		if err := s.scheduledConfirmationRepo.Create(ctx, scheduledConfirmation); err != nil {
			log.Printf("❌ Failed to create scheduled confirmation for property %s: %v", property.Reference, err)
			response.SkippedCount++
			response.SkippedReasons = append(response.SkippedReasons, fmt.Sprintf("Property %s: failed to save", property.Reference))
			continue
		}

		response.ScheduledCount++
		response.ScheduledConfirmIDs = append(response.ScheduledConfirmIDs, scheduledConfirmation.ID)

		log.Printf("✅ Scheduled confirmation for property %s (owner: %s, expires: %s)",
			property.Reference, property.OwnerID, expiresAt.Format("2006-01-02"))
	}

	log.Printf("📊 Scheduling complete: %d scheduled, %d skipped out of %d total properties",
		response.ScheduledCount, response.SkippedCount, response.TotalProperties)

	return response, nil
}

// ProcessPendingConfirmations processes all pending confirmations for today
// This can be called by a cron job daily to send scheduled confirmations
func (s *MonthlyConfirmationScheduler) ProcessPendingConfirmations(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	today := time.Now()
	log.Printf("🚀 Processing pending confirmations for tenant %s on %s", tenantID, today.Format("2006-01-02"))

	// Get all pending confirmations scheduled for today
	pendingConfirmations, err := s.scheduledConfirmationRepo.GetPendingForDate(ctx, tenantID, today)
	if err != nil {
		return fmt.Errorf("failed to get pending confirmations: %w", err)
	}

	log.Printf("📋 Found %d pending confirmations to process", len(pendingConfirmations))

	successCount := 0
	failCount := 0

	for _, sc := range pendingConfirmations {
		// TODO: When WhatsApp Business API is integrated, send message here
		// For now, just mark as "sent" (manual delivery)

		// Get owner info for logging
		owner, err := s.ownerRepo.Get(ctx, tenantID, sc.OwnerID)
		if err != nil {
			log.Printf("⚠️  Warning: could not get owner info for %s: %v", sc.OwnerID, err)
		}

		// Update status to sent
		now := time.Now()
		sc.Status = models.ScheduledConfirmationStatusSent
		sc.SentAt = &now
		sc.DeliveryStatus = "manual_delivery_required"

		if err := s.scheduledConfirmationRepo.Update(ctx, sc); err != nil {
			log.Printf("❌ Failed to update scheduled confirmation %s: %v", sc.ID, err)
			failCount++
			continue
		}

		successCount++
		ownerName := "Unknown"
		if owner != nil {
			ownerName = owner.Name
		}
		log.Printf("✅ Marked confirmation as ready for property %s (owner: %s)", sc.PropertyID, ownerName)
	}

	log.Printf("📊 Processing complete: %d successful, %d failed", successCount, failCount)
	return nil
}

// GetScheduledConfirmationsForTenant retrieves all scheduled confirmations for a tenant
func (s *MonthlyConfirmationScheduler) GetScheduledConfirmationsForTenant(ctx context.Context, tenantID string, status *models.ScheduledConfirmationStatus) ([]*models.ScheduledConfirmation, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant_id is required")
	}

	// Get all scheduled confirmations for the tenant
	confirmations, err := s.scheduledConfirmationRepo.ListByTenant(ctx, tenantID, status, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to list scheduled confirmations: %w", err)
	}

	return confirmations, nil
}

// GetScheduledConfirmationsForBroker retrieves scheduled confirmations for a broker to send manually
func (s *MonthlyConfirmationScheduler) GetScheduledConfirmationsForBroker(ctx context.Context, tenantID, brokerID string, status *models.ScheduledConfirmationStatus) ([]*models.ScheduledConfirmation, error) {
	if tenantID == "" || brokerID == "" {
		return nil, fmt.Errorf("tenant_id and broker_id are required")
	}

	// Get all scheduled confirmations for the tenant
	allConfirmations, err := s.scheduledConfirmationRepo.ListByTenant(ctx, tenantID, status, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to list scheduled confirmations: %w", err)
	}

	// Filter by broker
	var brokerConfirmations []*models.ScheduledConfirmation
	for _, sc := range allConfirmations {
		if sc.BrokerID == brokerID {
			brokerConfirmations = append(brokerConfirmations, sc)
		}
	}

	return brokerConfirmations, nil
}
