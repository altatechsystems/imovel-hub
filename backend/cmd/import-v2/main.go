package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/altatech/ecosistema-imob/backend/internal/adapters/union"
	"github.com/altatech/ecosistema-imob/backend/internal/services"
	"google.golang.org/api/option"
)

func main() {
	// Parse flags
	xmlFile := flag.String("xml", "", "Path to Union XML file (required)")
	xlsFile := flag.String("xls", "", "Path to Union XLS file (optional)")
	tenantID := flag.String("tenant", "", "Tenant ID (required)")
	createdBy := flag.String("created-by", "system", "Created by (broker_id or 'system')")
	limit := flag.Int("limit", 0, "Limit number of properties to import (0 = no limit)")

	flag.Parse()

	// Validate required flags
	if *xmlFile == "" {
		log.Fatal("--xml flag is required")
	}
	if *tenantID == "" {
		log.Fatal("--tenant flag is required")
	}

	ctx := context.Background()

	// Initialize Firestore with named database "imob-dev"
	projectID := "ecosistema-imob-dev"
	databaseID := "imob-dev"

	sa := option.WithCredentialsFile("config/firebase-adminsdk.json")
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, sa)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	log.Printf("✅ Connected to Firestore database: %s (project: %s)", databaseID, projectID)

	// Initialize import service
	importService := services.NewImportService(client)

	// Create import batch
	batch, err := importService.CreateBatch(ctx, *tenantID, "union", *createdBy)
	if err != nil {
		log.Fatalf("Failed to create import batch: %v", err)
	}

	log.Printf("✅ Created import batch: %s", batch.ID)

	// Parse XML
	log.Printf("📄 Parsing XML file: %s", *xmlFile)
	xmlData, err := parseXMLFile(*xmlFile)
	if err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	log.Printf("📊 Found %d properties in XML", len(xmlData.Imoveis))
	batch.TotalXMLRecords = len(xmlData.Imoveis)

	// Parse XLS (optional)
	var xlsRecords []union.XLSRecord
	if *xlsFile != "" {
		log.Printf("📄 Parsing XLS file: %s", *xlsFile)
		xlsRecords, err = union.ParseXLS(*xlsFile)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to parse XLS: %v", err)
		} else {
			log.Printf("📊 Found %d owner records in XLS", len(xlsRecords))
		}
	}

	// Import properties
	count := 0
	for i, xmlImovel := range xmlData.Imoveis {
		if *limit > 0 && count >= *limit {
			log.Printf("⏹️  Reached limit of %d properties", *limit)
			break
		}

		// Find matching XLS record
		var xlsRecord *union.XLSRecord
		if len(xlsRecords) > 0 {
			xlsRecord = union.FindXLSRecordByCode(xlsRecords, &xmlImovel)
		}

		// Normalize to PropertyPayload
		payload := union.NormalizeProperty(&xmlImovel, xlsRecord, *tenantID)

		// Import property (handles deduplication, owner, listing, etc.)
		err := importService.ImportProperty(ctx, batch, payload)
		if err != nil {
			log.Printf("❌ Error importing property %s: %v", xmlImovel.Referencia, err)

			// Log error to batch
			_ = importService.LogError(ctx, batch, "import_failed", err.Error(), map[string]interface{}{
				"reference":    xmlImovel.Referencia,
				"external_id":  xmlImovel.Codigoimovel,
				"property_idx": i,
			})

			continue
		}

		count++

		// Progress indicator
		if count%10 == 0 {
			log.Printf("📥 Imported %d/%d properties...", count, len(xmlData.Imoveis))
		}
	}

	// Complete batch
	if err := importService.CompleteBatch(ctx, batch); err != nil {
		log.Fatalf("Failed to complete batch: %v", err)
	}

	// Print summary
	fmt.Println("\n" + "═══════════════════════════════════════════════════════")
	fmt.Println("📊 IMPORT SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Printf("Batch ID:                       %s\n", batch.ID)
	fmt.Printf("Tenant ID:                      %s\n", batch.TenantID)
	fmt.Printf("Source:                         %s\n", batch.Source)
	fmt.Printf("Status:                         %s\n", batch.Status)
	fmt.Println("───────────────────────────────────────────────────────")
	fmt.Printf("Total XML Records:              %d\n", batch.TotalXMLRecords)
	fmt.Printf("Properties Created:             %d\n", batch.TotalPropertiesCreated)
	fmt.Printf("Properties Matched (Existing):  %d\n", batch.TotalPropertiesMatchedExisting)
	fmt.Printf("Possible Duplicates:            %d\n", batch.TotalPossibleDuplicates)
	fmt.Println("───────────────────────────────────────────────────────")
	fmt.Printf("Owners Enriched from XLS:       %d\n", batch.TotalOwnersEnrichedFromXLS)
	fmt.Printf("Owners Placeholders:            %d\n", batch.TotalOwnersPlaceholders)
	fmt.Println("───────────────────────────────────────────────────────")
	fmt.Printf("Listings Created:               %d\n", batch.TotalListingsCreated)
	fmt.Printf("Photos Processed:               %d\n", batch.TotalPhotosProcessed)
	fmt.Println("───────────────────────────────────────────────────────")
	fmt.Printf("Total Errors:                   %d\n", batch.TotalErrors)
	fmt.Println("═══════════════════════════════════════════════════════")

	if batch.TotalErrors > 0 {
		fmt.Println("\n⚠️  Some errors occurred during import. Check import_errors collection.")
	} else {
		fmt.Println("\n✅ Import completed successfully!")
	}
}

// parseXMLFile parses XML file
func parseXMLFile(filePath string) (*union.XMLUnion, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return union.ParseXML(file)
}
