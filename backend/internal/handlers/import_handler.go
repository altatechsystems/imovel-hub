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
	tenantID, exists := c.Get(middleware.TenantIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil { // 50 MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form", "details": err.Error()})
		return
	}

	// Get XML file
	xmlFile, xmlHeader, err := c.Request.FormFile("xml")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XML file is required", "details": err.Error()})
		return
	}
	defer xmlFile.Close()

	// Get optional XLS file
	var xlsFile io.ReadCloser
	var xlsHeader *http.Header
	xlsFileMultipart, xlsHeaderMultipart, err := c.Request.FormFile("xls")
	if err == nil {
		xlsFile = xlsFileMultipart
		h := http.Header(xlsHeaderMultipart.Header)
		xlsHeader = &h
		defer xlsFile.Close()
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

	// Save XML file
	xmlPath := filepath.Join(tempDir, xmlHeader.Filename)
	if err := saveUploadedFile(xmlFile, xmlPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save XML file", "details": err.Error()})
		return
	}

	// Save XLS file if provided
	var xlsPath string
	if xlsFile != nil && xlsHeader != nil {
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

	// Parse XLS if provided
	var xlsRecords []union.XLSRecord
	if xlsPath != "" {
		xlsRecords, err = union.ParseXLS(xlsPath)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to parse XLS: %v", err)
			// Continue without XLS data
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

// GetImportStatus handles GET /api/v1/tenants/{tenantId}/import/{batchId}
func (h *ImportHandler) GetImportStatus(c *gin.Context) {
	batchID := c.Param("batchId")

	// TODO: Implement batch status retrieval from Firestore
	// For now, return a simple response
	c.JSON(http.StatusOK, gin.H{
		"batch_id": batchID,
		"status":   "Check Firestore console for batch status",
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
