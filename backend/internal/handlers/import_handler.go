package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/altatech/ecosistema-imob/backend/internal/adapters/union"
	"github.com/altatech/ecosistema-imob/backend/internal/middleware"
	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/services"
)

// ImportHandler handles import operations
type ImportHandler struct {
	importService *services.ImportService
}

// NewImportHandler creates a new import handler
func NewImportHandler(importService *services.ImportService) *ImportHandler {
	return &ImportHandler{
		importService: importService,
	}
}

// ImportRequest represents the import request body
type ImportRequest struct {
	Source    string `json:"source" binding:"required"`    // "union"
	CreatedBy string `json:"created_by" binding:"required"` // broker_id or "system"
}

// ImportResponse represents the import response
type ImportResponse struct {
	BatchID   string                 `json:"batch_id"`
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Summary   map[string]interface{} `json:"summary,omitempty"`
}

// ImportFromFiles handles POST /api/v1/tenants/{tenantId}/import
// Accepts multipart form with XML and optional XLS files
func (h *ImportHandler) ImportFromFiles(c *gin.Context) {
	// Get tenant ID from middleware
	log.Printf("🔍 Checking for TenantID in context with key: %s", string(middleware.TenantIDKey))
	tenantID, exists := c.Get(string(middleware.TenantIDKey))
	if !exists {
		log.Printf("❌ TenantID not found in context!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}
	log.Printf("✅ TenantID found: %v", tenantID)

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil { // 50 MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form", "details": err.Error()})
		return
	}

	// Get XML file (optional now)
	xmlFile, xmlHeader, err := c.Request.FormFile("xml")
	hasXML := err == nil
	if hasXML {
		defer xmlFile.Close()
	}

	// Get XLS file (optional)
	var xlsFile io.ReadCloser
	xlsFileMultipart, xlsHeaderMultipart, err := c.Request.FormFile("xls")
	hasXLS := err == nil
	if hasXLS {
		xlsFile = xlsFileMultipart
		defer xlsFile.Close()
	}

	// At least one file must be provided
	if !hasXML && !hasXLS {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one file (XML or XLS) is required"})
		return
	}

	// Get source and created_by from form
	source := c.PostForm("source")
	if source == "" {
		source = "union" // default
	}

	createdBy := c.PostForm("created_by")
	if createdBy == "" {
		createdBy = "system" // default
	}

	// Save uploaded files to temp directory
	tempDir := filepath.Join(os.TempDir(), "import-"+uuid.New().String())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory"})
		return
	}
	defer os.RemoveAll(tempDir) // cleanup

	// Save XML file if provided
	var xmlPath string
	if hasXML {
		xmlPath = filepath.Join(tempDir, xmlHeader.Filename)
		if err := saveUploadedFile(xmlFile, xmlPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save XML file", "details": err.Error()})
			return
		}
	}

	// Save XLS file if provided
	var xlsPath string
	if hasXLS {
		xlsPath = filepath.Join(tempDir, xlsHeaderMultipart.Filename)
		if err := saveUploadedFile(xlsFile, xlsPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save XLS file", "details": err.Error()})
			return
		}
	}

	// Process import asynchronously
	ctx := context.Background()
	batch, err := h.importService.CreateBatch(ctx, tenantID.(string), source, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create import batch", "details": err.Error()})
		return
	}

	// Start async import
	go h.processImport(ctx, batch, xmlPath, xlsPath)

	// Return batch ID immediately
	c.JSON(http.StatusAccepted, ImportResponse{
		BatchID: batch.ID,
		Status:  "processing",
		Message: fmt.Sprintf("Import started. Batch ID: %s", batch.ID),
	})
}

// processImport processes the import in background
func (h *ImportHandler) processImport(ctx context.Context, batch *models.ImportBatch, xmlPath, xlsPath string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Import panic recovered: %v", r)
			batch.Status = "failed"
			_ = h.importService.CompleteBatch(ctx, batch)
		}
	}()

	// Parse XLS if provided
	var xlsRecords []union.XLSRecord
	if xlsPath != "" {
		var err error
		xlsRecords, err = union.ParseXLS(xlsPath)
		if err != nil {
			log.Printf("❌ Failed to parse XLS: %v", err)
			_ = h.importService.LogError(ctx, batch, "xls_parse", err.Error(), nil)
			batch.Status = "failed"
			_ = h.importService.CompleteBatch(ctx, batch)
			return
		}
		log.Printf("✅ Parsed XLS with %d records", len(xlsRecords))
	}

	// If only XLS provided (no XML), process XLS-only mode
	if xmlPath == "" && xlsPath != "" {
		h.processXLSOnlyImport(ctx, batch, xlsRecords)
		return
	}

	// Parse XML
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		log.Printf("❌ Failed to open XML file: %v", err)
		_ = h.importService.LogError(ctx, batch, "xml_open", err.Error(), nil)
		batch.Status = "failed"
		_ = h.importService.CompleteBatch(ctx, batch)
		return
	}
	defer xmlFile.Close()

	xmlData, err := union.ParseXML(xmlFile)
	if err != nil {
		log.Printf("❌ Failed to parse XML: %v", err)
		_ = h.importService.LogError(ctx, batch, "xml_parse", err.Error(), nil)
		batch.Status = "failed"
		_ = h.importService.CompleteBatch(ctx, batch)
		return
	}

	batch.TotalXMLRecords = len(xmlData.Imoveis)

	// Debug: Log photo count in first property from XML
	if len(xmlData.Imoveis) > 0 {
		firstProp := xmlData.Imoveis[0]
		log.Printf("🔍 DEBUG XML: First property %s has %d Foto tags in raw XML", firstProp.Referencia, len(firstProp.Fotos))
		if len(firstProp.Fotos) > 0 {
			log.Printf("   First Foto: URL='%s', Principal=%d", firstProp.Fotos[0].URL, firstProp.Fotos[0].Principal)
		}
	}

	// Import properties
	for i, xmlImovel := range xmlData.Imoveis {
		// Find matching XLS record
		var xlsRecord *union.XLSRecord
		if len(xlsRecords) > 0 {
			xlsRecord = union.FindXLSRecordByCode(xlsRecords, &xmlImovel)
		}

		// Normalize to PropertyPayload
		payload := union.NormalizeProperty(&xmlImovel, xlsRecord, batch.TenantID)

		// Debug: Log photo count
		if i < 3 { // Only log first 3 properties to avoid spam
			log.Printf("🖼️  Property %s has %d photos in payload", xmlImovel.Referencia, len(payload.Photos))
			if len(payload.Photos) > 0 {
				log.Printf("   First photo URL: %s", payload.Photos[0])
			}
		}

		// Import property
		if err := h.importService.ImportProperty(ctx, batch, payload); err != nil {
			log.Printf("❌ Error importing property %s: %v", xmlImovel.Referencia, err)
			_ = h.importService.LogError(ctx, batch, "import_failed", err.Error(), map[string]interface{}{
				"reference":    xmlImovel.Referencia,
				"external_id":  xmlImovel.Codigoimovel,
				"property_idx": i,
			})
			continue
		}

		// Progress log every 50 properties
		if (i+1)%50 == 0 {
			log.Printf("📥 Import progress: %d/%d properties", i+1, len(xmlData.Imoveis))
		}
	}

	// Complete batch
	if err := h.importService.CompleteBatch(ctx, batch); err != nil {
		log.Printf("❌ Failed to complete batch: %v", err)
	} else {
		log.Printf("✅ Import batch %s completed: %d properties created, %d errors", batch.ID, batch.TotalPropertiesCreated, batch.TotalErrors)
	}
}

// processXLSOnlyImport processes XLS-only import (update owner data for existing properties)
func (h *ImportHandler) processXLSOnlyImport(ctx context.Context, batch *models.ImportBatch, xlsRecords []union.XLSRecord) {
	log.Printf("🔄 Processing XLS-only import: %d records to process", len(xlsRecords))

	batch.TotalXMLRecords = len(xlsRecords) // Use XLS count as total
	batch.Status = "processing"

	for i, xlsRecord := range xlsRecords {
		// Build owner payload from XLS
		ownerPayload := union.OwnerPayload{
			Name:            xlsRecord.Proprietario,
			Email:           xlsRecord.Email,
			Phone:           xlsRecord.CelularTelefone,
			EnrichedFromXLS: true,
		}

		// Determine owner status
		if ownerPayload.Email != "" && ownerPayload.Phone != "" {
			ownerPayload.OwnerStatus = models.OwnerStatusVerified
		} else if ownerPayload.Phone != "" || ownerPayload.Email != "" {
			ownerPayload.OwnerStatus = models.OwnerStatusPartial
		} else {
			ownerPayload.OwnerStatus = models.OwnerStatusIncomplete
		}

		// Find existing property by reference
		propertyRef := xlsRecord.Referencia
		if propertyRef == "" {
			log.Printf("⚠️  Skipping XLS record %d: no reference code", i)
			continue
		}

		// Query property by reference
		property, err := h.importService.FindPropertyByReference(ctx, batch.TenantID, propertyRef)
		if err != nil || property == nil {
			log.Printf("⚠️  Property not found for reference %s, skipping", propertyRef)
			batch.TotalErrors++
			continue
		}

		// Update owner if property has one
		if property.OwnerID != "" {
			if err := h.importService.UpdateOwnerFromXLS(ctx, property.OwnerID, ownerPayload, propertyRef); err != nil {
				log.Printf("❌ Failed to update owner for property %s: %v", propertyRef, err)
				batch.TotalErrors++
			} else {
				batch.TotalOwnersEnrichedFromXLS++
				batch.TotalPropertiesMatchedExisting++
				log.Printf("✅ Updated owner for property %s", propertyRef)
			}
		} else {
			log.Printf("⚠️  Property %s has no owner, skipping", propertyRef)
		}

		// Progress log every 50 records
		if (i+1)%50 == 0 {
			log.Printf("📥 XLS import progress: %d/%d records", i+1, len(xlsRecords))
		}
	}

	// Log final statistics BEFORE completing
	log.Printf("📊 XLS-only import statistics BEFORE CompleteBatch:")
	log.Printf("   Batch ID: %s", batch.ID)
	log.Printf("   Total records: %d", batch.TotalXMLRecords)
	log.Printf("   Owners updated: %d", batch.TotalOwnersEnrichedFromXLS)
	log.Printf("   Properties matched: %d", batch.TotalPropertiesMatchedExisting)
	log.Printf("   Errors: %d", batch.TotalErrors)
	log.Printf("   Status: %s", batch.Status)

	// Complete batch
	if err := h.importService.CompleteBatch(ctx, batch); err != nil {
		log.Printf("❌ Failed to complete batch: %v", err)
	} else {
		log.Printf("✅ XLS-only import batch %s completed: %d owners updated, %d errors",
			batch.ID, batch.TotalOwnersEnrichedFromXLS, batch.TotalErrors)
	}
}

// GetImportStatus handles GET /api/v1/admin/:tenant_id/import/batches/:batchId
func (h *ImportHandler) GetImportStatus(c *gin.Context) {
	ctx := context.Background()
	batchID := c.Param("batchId")

	// Get batch from Firestore
	batch, err := h.importService.GetBatch(ctx, batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Batch not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, batch)
}

// saveUploadedFile saves an uploaded file to disk
func saveUploadedFile(src io.Reader, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
