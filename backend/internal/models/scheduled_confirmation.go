package models

import "time"

// ScheduledConfirmation tracks monthly confirmation reminders sent to property owners
// Collection: /tenants/{tenantId}/scheduled_confirmations/{id}
type ScheduledConfirmation struct {
	ID       string `firestore:"-" json:"id"`
	TenantID string `firestore:"tenant_id" json:"tenant_id"`

	// Property and Owner info
	PropertyID string `firestore:"property_id" json:"property_id"` // ref Property
	OwnerID    string `firestore:"owner_id" json:"owner_id"`       // ref Owner
	BrokerID   string `firestore:"broker_id" json:"broker_id"`     // ref Broker (captador)

	// Confirmation token info
	TokenID          string `firestore:"token_id" json:"token_id"`                     // ref OwnerConfirmationToken
	ConfirmationURL  string `firestore:"confirmation_url" json:"confirmation_url"`     // Full URL with token
	ConfirmationLink string `firestore:"confirmation_link" json:"confirmation_link"`   // Short link (if applicable)

	// Scheduling info
	ScheduledFor time.Time                  `firestore:"scheduled_for" json:"scheduled_for"` // When it should be sent
	SentAt       *time.Time                 `firestore:"sent_at,omitempty" json:"sent_at,omitempty"`
	Status       ScheduledConfirmationStatus `firestore:"status" json:"status"` // pending, sent, failed, cancelled

	// Delivery info
	DeliveryMethod string `firestore:"delivery_method" json:"delivery_method"` // whatsapp, email, sms, manual
	DeliveryStatus string `firestore:"delivery_status,omitempty" json:"delivery_status,omitempty"`
	DeliveryError  string `firestore:"delivery_error,omitempty" json:"delivery_error,omitempty"`

	// Owner response tracking
	RespondedAt *time.Time `firestore:"responded_at,omitempty" json:"responded_at,omitempty"`
	Response    string     `firestore:"response,omitempty" json:"response,omitempty"` // available, unavailable, price_updated

	// Metadata
	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}

// ScheduledConfirmationStatus represents the status of a scheduled confirmation
type ScheduledConfirmationStatus string

const (
	ScheduledConfirmationStatusPending   ScheduledConfirmationStatus = "pending"
	ScheduledConfirmationStatusSent      ScheduledConfirmationStatus = "sent"
	ScheduledConfirmationStatusFailed    ScheduledConfirmationStatus = "failed"
	ScheduledConfirmationStatusCancelled ScheduledConfirmationStatus = "cancelled"
	ScheduledConfirmationStatusResponded ScheduledConfirmationStatus = "responded"
)
