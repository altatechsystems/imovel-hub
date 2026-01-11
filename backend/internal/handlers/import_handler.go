package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	// Note: cleanup is done in processImport goroutine, not here

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
	// Clean up temp directory after import completes
	if xmlPath != "" || xlsPath != "" {
		var tempDir string
		if xmlPath != "" {
			tempDir = filepath.Dir(xmlPath)
		} else {
			tempDir = filepath.Dir(xlsPath)
		}
		defer func() {
			log.Printf("🧹 Cleaning up temp directory: %s", tempDir)
			os.RemoveAll(tempDir)
		}()
	}

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
		log.Printf("🎯 Detected XLS-only import mode (no XML file)")
		log.Printf("🎯 XLS records count: %d", len(xlsRecords))
		if len(xlsRecords) > 0 {
			log.Printf("🎯 First record: Ref=%s, Owner=%s, Phone=%s, Email=%s",
				xlsRecords[0].Referencia, xlsRecords[0].Proprietario,
				xlsRecords[0].CelularTelefone, xlsRecords[0].Email)
		}
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

	// Import properties in batches with concurrency control
	const batchSize = 50  // Process 50 properties at a time
	const maxWorkers = 3  // Limit concurrent goroutines to reduce memory usage

	totalProperties := len(xmlData.Imoveis)
	semaphore := make(chan struct{}, maxWorkers)

	for i := 0; i < totalProperties; i += batchSize {
		end := i + batchSize
		if end > totalProperties {
			end = totalProperties
		}

		batchProperties := xmlData.Imoveis[i:end]
		log.Printf("📦 Processing batch %d-%d of %d properties", i+1, end, totalProperties)

		// Process each property in this batch
		for idx, xmlImovel := range batchProperties {
			globalIdx := i + idx

			// Acquire semaphore (block if maxWorkers goroutines are running)
			semaphore <- struct{}{}

			// Process property (synchronously for now to avoid overwhelming Firestore)
			func(imovel union.XMLImovel, propertyIdx int) {
				defer func() { <-semaphore }() // Release semaphore

				// Find matching XLS record
				var xlsRecord *union.XLSRecord
				if len(xlsRecords) > 0 {
					xlsRecord = union.FindXLSRecordByCode(xlsRecords, &imovel)
				}

				// Normalize to PropertyPayload
				payload := union.NormalizeProperty(&imovel, xlsRecord, batch.TenantID)

				// Debug: Log photo count
				if propertyIdx < 3 {
					log.Printf("🖼️  Property %s has %d photos in payload", imovel.Referencia, len(payload.Photos))
					if len(payload.Photos) > 0 {
						log.Printf("   First photo URL: %s", payload.Photos[0])
					}
				}

				// Import property
				if err := h.importService.ImportProperty(ctx, batch, payload); err != nil {
					log.Printf("❌ Error importing property %s: %v", imovel.Referencia, err)
					_ = h.importService.LogError(ctx, batch, "import_failed", err.Error(), map[string]interface{}{
						"reference":    imovel.Referencia,
						"external_id":  imovel.Codigoimovel,
						"property_idx": propertyIdx,
					})
				}
			}(xmlImovel, globalIdx)
		}

		// Wait for all goroutines in this batch to complete
		for j := 0; j < maxWorkers; j++ {
			semaphore <- struct{}{}
		}
		for j := 0; j < maxWorkers; j++ {
			<-semaphore
		}

		log.Printf("✅ Batch complete: %d/%d properties processed", end, totalProperties)

		// Optional: Add small delay between batches to allow GC to run
		if end < totalProperties {
			log.Printf("⏸️  Pausing 2s between batches to allow memory cleanup...")
			time.Sleep(2 * time.Second)
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

	// Save batch immediately to update status and total count in Firestore
	if _, err := h.importService.GetDB().Collection("import_batches").Doc(batch.ID).Set(ctx, batch); err != nil {
		log.Printf("❌ Failed to save batch initial state: %v", err)
	} else {
		log.Printf("✅ Saved initial batch state: TotalXMLRecords=%d", batch.TotalXMLRecords)
	}

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

		// Log detailed info for CH00038
		if propertyRef == "CH00038" {
			log.Printf("🔍 Found CH00038 in XLS! Owner=%s, Phone=%s, Email=%s",
				xlsRecord.Proprietario, xlsRecord.CelularTelefone, xlsRecord.Email)
		}

		// Query property by reference
		property, err := h.importService.FindPropertyByReference(ctx, batch.TenantID, propertyRef)
		if err != nil {
			log.Printf("⚠️  Error finding property %s: %v", propertyRef, err)
			batch.TotalErrors++
			continue
		}
		if property == nil {
			log.Printf("⚠️  Property not found for reference %s (tenant: %s)", propertyRef, batch.TenantID)
			batch.TotalErrors++
			continue
		}

		// Log found property
		if propertyRef == "CH00038" {
			log.Printf("🔍 Found property CH00038 in DB! OwnerID=%s", property.OwnerID)
		}

		// Update owner if property has one
		if property.OwnerID != "" {
			if err := h.importService.UpdateOwnerFromXLS(ctx, batch.TenantID, property.OwnerID, ownerPayload, propertyRef); err != nil {
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

// GetBatchErrors handles GET /api/v1/admin/:tenant_id/import/batches/:batchId/errors
func (h *ImportHandler) GetBatchErrors(c *gin.Context) {
	ctx := context.Background()
	batchID := c.Param("batchId")

	// Get errors from Firestore
	errors, err := h.importService.GetBatchErrors(ctx, batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch errors",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"errors":  errors,
		"count":   len(errors),
	})
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
