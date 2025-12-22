package services

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"google.golang.org/api/iterator"
)

// DeduplicationService handles property deduplication logic
type DeduplicationService struct {
	db *firestore.Client
}

// NewDeduplicationService creates a new deduplication service
func NewDeduplicationService(db *firestore.Client) *DeduplicationService {
	return &DeduplicationService{db: db}
}

// DeduplicationResult contains deduplication check results
type DeduplicationResult struct {
	IsDuplicate       bool
	ExistingProperty  *models.Property
	MatchType         string // "external_id", "fingerprint", "none"
	PossibleDuplicate bool   // true if fingerprint matches but external_id doesn't
}

// CheckDuplicate checks if property already exists
// Priority: 1) external_source + external_id (strong match)
//           2) fingerprint (heuristic match)
func (s *DeduplicationService) CheckDuplicate(ctx context.Context, property *models.Property) (*DeduplicationResult, error) {
	result := &DeduplicationResult{
		IsDuplicate:       false,
		PossibleDuplicate: false,
		MatchType:         "none",
	}

	// 1. Check by external_source + external_id (strongest match)
	if property.ExternalSource != "" && property.ExternalID != "" {
		existing, err := s.findByExternalID(ctx, property.TenantID, property.ExternalSource, property.ExternalID)
		if err != nil {
			return nil, fmt.Errorf("failed to check external_id: %w", err)
		}

		if existing != nil {
			result.IsDuplicate = true
			result.ExistingProperty = existing
			result.MatchType = "external_id"
			return result, nil
		}
	}

	// 2. Check by fingerprint (heuristic match)
	if property.Fingerprint != "" {
		existing, err := s.findByFingerprint(ctx, property.TenantID, property.Fingerprint)
		if err != nil {
			return nil, fmt.Errorf("failed to check fingerprint: %w", err)
		}

		if existing != nil {
			// Fingerprint matches but external_id didn't
			// This is a POSSIBLE duplicate (needs manual review)
			result.PossibleDuplicate = true
			result.ExistingProperty = existing
			result.MatchType = "fingerprint"
			return result, nil
		}
	}

	return result, nil
}

// findByExternalID finds property by external source and ID
func (s *DeduplicationService) findByExternalID(ctx context.Context, tenantID, externalSource, externalID string) (*models.Property, error) {
	iter := s.db.Collection("properties").
		Where("tenant_id", "==", tenantID).
		Where("external_source", "==", externalSource).
		Where("external_id", "==", externalID).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var property models.Property
	if err := doc.DataTo(&property); err != nil {
		return nil, err
	}

	property.ID = doc.Ref.ID
	return &property, nil
}

// findByFingerprint finds property by fingerprint hash
func (s *DeduplicationService) findByFingerprint(ctx context.Context, tenantID, fingerprint string) (*models.Property, error) {
	iter := s.db.Collection("properties").
		Where("tenant_id", "==", tenantID).
		Where("fingerprint", "==", fingerprint).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var property models.Property
	if err := doc.DataTo(&property); err != nil {
		return nil, err
	}

	property.ID = doc.Ref.ID
	return &property, nil
}

// MarkAsPossibleDuplicate marks property as possible duplicate
func (s *DeduplicationService) MarkAsPossibleDuplicate(ctx context.Context, propertyID string, relatedPropertyID string) error {
	_, err := s.db.Collection("properties").Doc(propertyID).Update(ctx, []firestore.Update{
		{Path: "possible_duplicate", Value: true},
		{Path: "related_property_id", Value: relatedPropertyID},
	})
	return err
}
