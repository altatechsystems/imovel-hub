package models

import "time"

// ImportBatch represents a batch import operation
type ImportBatch struct {
	ID       string `firestore:"-" json:"id"`
	TenantID string `firestore:"tenant_id" json:"tenant_id"`
	Source   string `firestore:"source" json:"source"` // "union"
	Status   string `firestore:"status" json:"status"` // processing, completed, failed

	// Summary counters
	TotalXMLRecords                int `firestore:"total_xml_records" json:"total_xml_records"`
	TotalPropertiesCreated         int `firestore:"total_properties_created" json:"total_properties_created"`
	TotalPropertiesMatchedExisting int `firestore:"total_properties_matched_existing" json:"total_properties_matched_existing"`
	TotalPossibleDuplicates        int `firestore:"total_possible_duplicates" json:"total_possible_duplicates"`
	TotalOwnersPlaceholders        int `firestore:"total_owners_placeholders" json:"total_owners_placeholders"`
	TotalOwnersEnrichedFromXLS     int `firestore:"total_owners_enriched_from_xls" json:"total_owners_enriched_from_xls"`
	TotalListingsCreated           int `firestore:"total_listings_created" json:"total_listings_created"`
	TotalPhotosProcessed           int `firestore:"total_photos_processed" json:"total_photos_processed"`
	TotalErrors                    int `firestore:"total_errors" json:"total_errors"`

	// Metadata
	StartedAt   time.Time  `firestore:"started_at" json:"started_at"`
	CompletedAt *time.Time `firestore:"completed_at,omitempty" json:"completed_at,omitempty"`
	CreatedBy   string     `firestore:"created_by" json:"created_by"` // broker_id or system
}

// ImportError represents an error that occurred during import
type ImportError struct {
	ID           string                 `firestore:"-" json:"id"`
	BatchID      string                 `firestore:"batch_id" json:"batch_id"`
	TenantID     string                 `firestore:"tenant_id" json:"tenant_id"`
	ErrorType    string                 `firestore:"error_type" json:"error_type"` // xml_parse, photo_download, deduplication, etc.
	ErrorMessage string                 `firestore:"error_message" json:"error_message"`
	RecordData   map[string]interface{} `firestore:"record_data" json:"record_data"` // dados do registro com problema
	Timestamp    time.Time              `firestore:"timestamp" json:"timestamp"`
}
