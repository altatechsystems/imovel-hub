package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
)

// ScheduledConfirmationRepository handles Firestore operations for scheduled confirmations
type ScheduledConfirmationRepository struct {
	*BaseRepository
}

// NewScheduledConfirmationRepository creates a new scheduled confirmation repository
func NewScheduledConfirmationRepository(client *firestore.Client) *ScheduledConfirmationRepository {
	return &ScheduledConfirmationRepository{
		BaseRepository: NewBaseRepository(client),
	}
}

// Create creates a new scheduled confirmation
func (r *ScheduledConfirmationRepository) Create(ctx context.Context, sc *models.ScheduledConfirmation) error {
	if sc.TenantID == "" {
		return fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}

	if sc.ID == "" {
		sc.ID = r.GenerateID("scheduled_confirmations")
	}

	now := time.Now()
	sc.CreatedAt = now
	sc.UpdatedAt = now

	if err := r.CreateDocument(ctx, "scheduled_confirmations", sc.ID, sc); err != nil {
		return fmt.Errorf("failed to create scheduled confirmation: %w", err)
	}

	return nil
}

// Get retrieves a scheduled confirmation by ID
func (r *ScheduledConfirmationRepository) Get(ctx context.Context, tenantID, id string) (*models.ScheduledConfirmation, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}

	var sc models.ScheduledConfirmation
	if err := r.GetDocument(ctx, "scheduled_confirmations", id, &sc); err != nil {
		return nil, err
	}

	// Verify tenant ownership
	if sc.TenantID != tenantID {
		return nil, ErrNotFound
	}

	sc.ID = id
	return &sc, nil
}

// Update updates a scheduled confirmation
func (r *ScheduledConfirmationRepository) Update(ctx context.Context, sc *models.ScheduledConfirmation) error {
	if sc.TenantID == "" || sc.ID == "" {
		return fmt.Errorf("%w: tenant_id and id are required", ErrInvalidInput)
	}

	sc.UpdatedAt = time.Now()

	updates := []firestore.Update{
		{Path: "status", Value: sc.Status},
		{Path: "sent_at", Value: sc.SentAt},
		{Path: "delivery_status", Value: sc.DeliveryStatus},
		{Path: "delivery_error", Value: sc.DeliveryError},
		{Path: "responded_at", Value: sc.RespondedAt},
		{Path: "response", Value: sc.Response},
		{Path: "updated_at", Value: sc.UpdatedAt},
	}

	if err := r.UpdateDocument(ctx, "scheduled_confirmations", sc.ID, updates); err != nil {
		return fmt.Errorf("failed to update scheduled confirmation: %w", err)
	}

	return nil
}

// GetPendingForDate retrieves all pending scheduled confirmations for a specific date
func (r *ScheduledConfirmationRepository) GetPendingForDate(ctx context.Context, tenantID string, targetDate time.Time) ([]*models.ScheduledConfirmation, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}

	// Query for pending confirmations scheduled for the target date
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := r.Client().Collection("scheduled_confirmations").
		Where("tenant_id", "==", tenantID).
		Where("status", "==", string(models.ScheduledConfirmationStatusPending)).
		Where("scheduled_for", ">=", startOfDay).
		Where("scheduled_for", "<", endOfDay).
		OrderBy("scheduled_for", firestore.Asc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var results []*models.ScheduledConfirmation
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate scheduled confirmations: %w", err)
		}

		var sc models.ScheduledConfirmation
		if err := doc.DataTo(&sc); err != nil {
			return nil, fmt.Errorf("failed to decode scheduled confirmation: %w", err)
		}

		sc.ID = doc.Ref.ID
		results = append(results, &sc)
	}

	return results, nil
}

// GetByPropertyAndMonth retrieves scheduled confirmations for a property in a specific month
func (r *ScheduledConfirmationRepository) GetByPropertyAndMonth(ctx context.Context, tenantID, propertyID string, year int, month time.Month) ([]*models.ScheduledConfirmation, error) {
	if tenantID == "" || propertyID == "" {
		return nil, fmt.Errorf("%w: tenant_id and property_id are required", ErrInvalidInput)
	}

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	query := r.Client().Collection("scheduled_confirmations").
		Where("tenant_id", "==", tenantID).
		Where("property_id", "==", propertyID).
		Where("scheduled_for", ">=", startOfMonth).
		Where("scheduled_for", "<", endOfMonth).
		OrderBy("scheduled_for", firestore.Desc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var results []*models.ScheduledConfirmation
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate scheduled confirmations: %w", err)
		}

		var sc models.ScheduledConfirmation
		if err := doc.DataTo(&sc); err != nil {
			return nil, fmt.Errorf("failed to decode scheduled confirmation: %w", err)
		}

		sc.ID = doc.Ref.ID
		results = append(results, &sc)
	}

	return results, nil
}

// ListByTenant retrieves all scheduled confirmations for a tenant with optional status filter
func (r *ScheduledConfirmationRepository) ListByTenant(ctx context.Context, tenantID string, status *models.ScheduledConfirmationStatus, limit int) ([]*models.ScheduledConfirmation, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant_id is required", ErrInvalidInput)
	}

	query := r.Client().Collection("scheduled_confirmations").
		Where("tenant_id", "==", tenantID)

	if status != nil {
		query = query.Where("status", "==", string(*status))
	}

	// Remove OrderBy to avoid needing composite index
	// We'll sort in memory instead

	if limit > 0 {
		query = query.Limit(limit)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var results []*models.ScheduledConfirmation
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate scheduled confirmations: %w", err)
		}

		var sc models.ScheduledConfirmation
		if err := doc.DataTo(&sc); err != nil {
			return nil, fmt.Errorf("failed to decode scheduled confirmation: %w", err)
		}

		sc.ID = doc.Ref.ID
		results = append(results, &sc)
	}

	return results, nil
}
