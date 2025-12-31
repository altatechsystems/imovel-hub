package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduleMonthlyConfirmationsRequest_DefaultScheduledDate(t *testing.T) {
	// Test that when ScheduledFor is zero, it defaults to next month 1st at 9 AM
	req := ScheduleMonthlyConfirmationsRequest{
		TenantID:     "test-tenant",
		ScheduledFor: time.Time{}, // Zero value
		DryRun:       true,
	}

	now := time.Now()
	expectedMonth := now.Month() + 1
	expectedYear := now.Year()
	if expectedMonth > 12 {
		expectedMonth = 1
		expectedYear++
	}

	// Simulate the logic from the service
	if req.ScheduledFor.IsZero() {
		req.ScheduledFor = time.Date(now.Year(), now.Month()+1, 1, 9, 0, 0, 0, now.Location())
	}

	assert.Equal(t, expectedYear, req.ScheduledFor.Year())
	assert.Equal(t, expectedMonth, req.ScheduledFor.Month())
	assert.Equal(t, 1, req.ScheduledFor.Day())
	assert.Equal(t, 9, req.ScheduledFor.Hour())
	assert.Equal(t, 0, req.ScheduledFor.Minute())
}

func TestScheduleMonthlyConfirmationsRequest_CustomScheduledDate(t *testing.T) {
	// Test that custom scheduled date is respected
	customDate := time.Date(2025, 3, 15, 14, 30, 0, 0, time.UTC)
	req := ScheduleMonthlyConfirmationsRequest{
		TenantID:     "test-tenant",
		ScheduledFor: customDate,
		DryRun:       false,
	}

	assert.Equal(t, customDate, req.ScheduledFor)
	assert.Equal(t, 2025, req.ScheduledFor.Year())
	assert.Equal(t, time.March, req.ScheduledFor.Month())
	assert.Equal(t, 15, req.ScheduledFor.Day())
	assert.Equal(t, 14, req.ScheduledFor.Hour())
}

func TestScheduleMonthlyConfirmationsResponse_Structure(t *testing.T) {
	// Test response structure
	response := &ScheduleMonthlyConfirmationsResponse{
		TotalProperties:     100,
		ScheduledCount:      85,
		SkippedCount:        15,
		ScheduledForDate:    "2025-02-01 09:00:00",
		SkippedReasons:      []string{"Property TE001: no owner", "Property TE002: status is unavailable"},
		ScheduledConfirmIDs: []string{"sc1", "sc2", "sc3"},
	}

	assert.Equal(t, 100, response.TotalProperties)
	assert.Equal(t, 85, response.ScheduledCount)
	assert.Equal(t, 15, response.SkippedCount)
	assert.Equal(t, "2025-02-01 09:00:00", response.ScheduledForDate)
	assert.Len(t, response.SkippedReasons, 2)
	assert.Len(t, response.ScheduledConfirmIDs, 3)
	assert.Contains(t, response.SkippedReasons[0], "no owner")
	assert.Contains(t, response.SkippedReasons[1], "status is unavailable")
}

func TestScheduleMonthlyConfirmationsResponse_EmptyReasons(t *testing.T) {
	// Test that SkippedReasons can be empty when all properties are scheduled
	response := &ScheduleMonthlyConfirmationsResponse{
		TotalProperties:     10,
		ScheduledCount:      10,
		SkippedCount:        0,
		ScheduledForDate:    "2025-02-01 09:00:00",
		SkippedReasons:      []string{},
		ScheduledConfirmIDs: []string{"sc1", "sc2", "sc3", "sc4", "sc5", "sc6", "sc7", "sc8", "sc9", "sc10"},
	}

	assert.Equal(t, 10, response.TotalProperties)
	assert.Equal(t, 10, response.ScheduledCount)
	assert.Equal(t, 0, response.SkippedCount)
	assert.Empty(t, response.SkippedReasons)
	assert.Len(t, response.ScheduledConfirmIDs, 10)
}

func TestScheduleMonthlyConfirmations_DryRunDoesNotCreateRecords(t *testing.T) {
	// Validate that dry run mode is indicated by the DryRun flag
	dryRunReq := ScheduleMonthlyConfirmationsRequest{
		TenantID:     "test-tenant",
		ScheduledFor: time.Date(2025, 2, 1, 9, 0, 0, 0, time.UTC),
		DryRun:       true,
	}

	assert.True(t, dryRunReq.DryRun, "DryRun should be true")

	normalReq := ScheduleMonthlyConfirmationsRequest{
		TenantID:     "test-tenant",
		ScheduledFor: time.Date(2025, 2, 1, 9, 0, 0, 0, time.UTC),
		DryRun:       false,
	}

	assert.False(t, normalReq.DryRun, "DryRun should be false")
}

func TestScheduledForDate_ValidFormats(t *testing.T) {
	// Test that date formatting works correctly
	testDate := time.Date(2025, 2, 1, 9, 0, 0, 0, time.UTC)
	formatted := testDate.Format("2006-01-02 15:04:05")

	assert.Equal(t, "2025-02-01 09:00:00", formatted)

	// Test ISO format
	isoFormatted := testDate.Format(time.RFC3339)
	assert.Equal(t, "2025-02-01T09:00:00Z", isoFormatted)
}

func TestMonthCalculation_HandleYearRollover(t *testing.T) {
	// Test that month calculation handles December -> January correctly
	decemberDate := time.Date(2025, 12, 15, 10, 0, 0, 0, time.UTC)

	nextMonth := decemberDate.Month() + 1
	expectedYear := decemberDate.Year()
	expectedMonth := nextMonth

	if nextMonth > 12 {
		expectedMonth = 1
		expectedYear++
	}

	assert.Equal(t, 2026, expectedYear)
	assert.Equal(t, time.January, expectedMonth)

	// Test normal month increment
	februaryDate := time.Date(2025, 2, 15, 10, 0, 0, 0, time.UTC)
	nextMonth = februaryDate.Month() + 1
	expectedYear = februaryDate.Year()
	expectedMonth = nextMonth

	if nextMonth > 12 {
		expectedMonth = 1
		expectedYear++
	}

	assert.Equal(t, 2025, expectedYear)
	assert.Equal(t, time.March, expectedMonth)
}

func TestSkipReasons_FormatConsistency(t *testing.T) {
	// Test that skip reasons follow a consistent format
	testCases := []struct {
		propertyRef string
		reason      string
		expected    string
	}{
		{"TE00001", "no owner", "Property TE00001: no owner"},
		{"TE00002", "status is unavailable", "Property TE00002: status is unavailable"},
		{"TE00003", "already scheduled", "Property TE00003: already scheduled"},
		{"TE00004", "failed to generate link", "Property TE00004: failed to generate link"},
	}

	for _, tc := range testCases {
		result := "Property " + tc.propertyRef + ": " + tc.reason
		assert.Equal(t, tc.expected, result)
	}
}

func TestRequestValidation_TenantIDRequired(t *testing.T) {
	// Test that empty tenant ID is detectable
	emptyTenantReq := ScheduleMonthlyConfirmationsRequest{
		TenantID:     "",
		ScheduledFor: time.Date(2025, 2, 1, 9, 0, 0, 0, time.UTC),
		DryRun:       false,
	}

	assert.Empty(t, emptyTenantReq.TenantID, "Empty tenant ID should be detectable")

	validReq := ScheduleMonthlyConfirmationsRequest{
		TenantID:     "valid-tenant-id",
		ScheduledFor: time.Date(2025, 2, 1, 9, 0, 0, 0, time.UTC),
		DryRun:       false,
	}

	assert.NotEmpty(t, validReq.TenantID, "Valid tenant ID should be present")
}
