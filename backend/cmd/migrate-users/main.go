package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
)

// MigrationStats tracks migration statistics
type MigrationStats struct {
	TotalBrokers       int
	RealBrokers        int // With valid CRECI
	AdminUsers         int // Without CRECI or with admin/manager role
	Migrated           int
	Errors             int
	Skipped            int
	ErrorDetails       []string
}

func main() {
	// Command line flags
	dryRun := flag.Bool("dry-run", true, "Run in dry-run mode (no actual changes)")
	csvReport := flag.String("csv", "migration-report.csv", "CSV report file")
	projectID := flag.String("project-id", "", "Firebase project ID (or set FIREBASE_PROJECT_ID env var)")
	flag.Parse()

	log.Println("===========================================")
	log.Println("MIGRATION: Separate Brokers from Users")
	log.Println("===========================================")
	if *dryRun {
		log.Println("MODE: DRY-RUN (no changes will be made)")
	} else {
		log.Println("MODE: EXECUTION (changes will be applied)")
	}
	log.Println("===========================================\n")

	// Get project ID
	if *projectID == "" {
		*projectID = os.Getenv("FIREBASE_PROJECT_ID")
	}
	if *projectID == "" {
		*projectID = "ecosistema-imob-dev" // Default project ID
	}

	log.Printf("Using Firebase Project ID: %s\n\n", *projectID)

	// Initialize Firebase
	ctx := context.Background()

	// Check for service account key
	credPath := "config/firebase-adminsdk.json"
	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		log.Fatalf("Firebase service account key not found at %s", credPath)
	}

	opt := option.WithCredentialsFile(credPath)

	// Connect to the named Firestore database "imob-dev" (where data is actually stored)
	log.Println("Connecting to Firestore database: imob-dev")
	client, err := firestore.NewClientWithDatabase(ctx, *projectID, "imob-dev", opt)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
	}
	defer client.Close()
	log.Println("✅ Connected to Firestore database: imob-dev")

	// Initialize stats
	stats := &MigrationStats{
		ErrorDetails: make([]string, 0),
	}

	// Get all tenants
	tenants, err := getAllTenants(ctx, client)
	if err != nil {
		log.Fatalf("Error getting tenants: %v", err)
	}

	log.Printf("Found %d tenant(s) to process\n\n", len(tenants))

	// Open CSV report
	csvFile, err := os.Create(*csvReport)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	csvWriter.Write([]string{
		"Tenant ID",
		"Tenant Name",
		"Broker ID",
		"Name",
		"Email",
		"CRECI",
		"Role",
		"Type",
		"Action",
		"Status",
		"Notes",
	})

	// Process each tenant
	for _, tenant := range tenants {
		log.Printf("Processing tenant: %s (%s)\n", tenant.Name, tenant.ID)
		processTenant(ctx, client, tenant, stats, csvWriter, *dryRun)
		log.Println()
	}

	// Print final stats
	printStats(stats)

	log.Printf("\nMigration report saved to: %s\n", *csvReport)
}

func getAllTenants(ctx context.Context, client *firestore.Client) ([]*models.Tenant, error) {
	var tenants []*models.Tenant

	iter := client.Collection("tenants").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var tenant models.Tenant
		if err := doc.DataTo(&tenant); err != nil {
			log.Printf("Warning: Failed to parse tenant %s: %v\n", doc.Ref.ID, err)
			continue
		}

		tenant.ID = doc.Ref.ID
		tenants = append(tenants, &tenant)
	}

	return tenants, nil
}

func processTenant(ctx context.Context, client *firestore.Client, tenant *models.Tenant, stats *MigrationStats, csvWriter *csv.Writer, dryRun bool) {
	// Get all brokers from this tenant
	brokersRef := client.Collection("tenants").Doc(tenant.ID).Collection("brokers")
	iter := brokersRef.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("  Error reading broker: %v\n", err)
			stats.Errors++
			continue
		}

		var broker models.Broker
		if err := doc.DataTo(&broker); err != nil {
			log.Printf("  Error parsing broker %s: %v\n", doc.Ref.ID, err)
			stats.Errors++
			continue
		}

		broker.ID = doc.Ref.ID
		stats.TotalBrokers++

		// Analyze broker
		brokerType, action, shouldMigrate := analyzeBroker(&broker)

		// Log analysis
		notes := ""
		status := "OK"

		if shouldMigrate {
			if dryRun {
				log.Printf("  [DRY-RUN] Would migrate: %s (%s) - %s -> User\n", broker.Name, broker.Email, brokerType)
				stats.AdminUsers++
			} else {
				// Perform actual migration
				err := migrateToUser(ctx, client, tenant.ID, &broker)
				if err != nil {
					log.Printf("  [ERROR] Failed to migrate %s: %v\n", broker.Name, err)
					stats.Errors++
					stats.ErrorDetails = append(stats.ErrorDetails, fmt.Sprintf("%s (%s): %v", broker.Name, broker.Email, err))
					status = "ERROR"
					notes = err.Error()
				} else {
					log.Printf("  [MIGRATED] %s (%s) - %s -> User\n", broker.Name, broker.Email, brokerType)
					stats.Migrated++
					stats.AdminUsers++
					status = "MIGRATED"
				}
			}
		} else {
			log.Printf("  [KEEP] %s (%s) - %s (CRECI: %s)\n", broker.Name, broker.Email, brokerType, broker.CRECI)
			stats.RealBrokers++
			stats.Skipped++
		}

		// Write to CSV
		csvWriter.Write([]string{
			tenant.ID,
			tenant.Name,
			broker.ID,
			broker.Name,
			broker.Email,
			broker.CRECI,
			broker.Role,
			brokerType,
			action,
			status,
			notes,
		})
	}
}

// analyzeBroker determines if a broker should be migrated to user
func analyzeBroker(broker *models.Broker) (brokerType string, action string, shouldMigrate bool) {
	// Check CRECI
	hasCRECI := broker.CRECI != "" &&
		broker.CRECI != "-" &&
		broker.CRECI != "PENDENTE" &&
		!strings.Contains(strings.ToLower(broker.CRECI), "pending") &&
		!strings.Contains(strings.ToLower(broker.CRECI), "n/a") &&
		len(broker.CRECI) > 3

	// Check role
	hasAdminRole := broker.Role == "admin" || broker.Role == "manager"
	hasBrokerRole := broker.Role == "broker" || broker.Role == "broker_admin"

	// Decision logic
	if !hasCRECI && hasAdminRole {
		return "Admin User (No CRECI)", "Migrate to /users", true
	}

	if !hasCRECI && !hasBrokerRole {
		return "Admin User (No CRECI, No Broker Role)", "Migrate to /users", true
	}

	if !hasCRECI {
		return "Invalid Broker (No CRECI)", "Migrate to /users", true
	}

	if hasCRECI && hasAdminRole {
		return "Broker Admin (Has CRECI)", "Keep in /brokers, update role to broker_admin", false
	}

	if hasCRECI {
		return "Real Broker", "Keep in /brokers", false
	}

	return "Unknown", "Manual Review Required", false
}

// migrateToUser migrates a broker to the users collection
func migrateToUser(ctx context.Context, client *firestore.Client, tenantID string, broker *models.Broker) error {
	// Create user from broker
	user := &models.User{
		ID:          broker.ID, // Keep same ID for traceability
		TenantID:    tenantID,
		FirebaseUID: broker.FirebaseUID,
		Name:        broker.Name,
		Email:       broker.Email,
		Phone:       broker.Phone,
		Document:    broker.Document,
		DocumentType: broker.DocumentType,
		Role:        determineUserRole(broker.Role),
		IsActive:    broker.IsActive,
		PhotoURL:    broker.PhotoURL,
		Permissions: []string{}, // Start with no permissions, add as needed
		CreatedAt:   broker.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// Create user in /users collection
	userRef := client.Collection("tenants").Doc(tenantID).Collection("users").Doc(user.ID)
	_, err := userRef.Set(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Delete from /brokers collection
	brokerRef := client.Collection("tenants").Doc(tenantID).Collection("brokers").Doc(broker.ID)
	_, err = brokerRef.Delete(ctx)
	if err != nil {
		// If deletion fails, try to rollback user creation
		userRef.Delete(ctx)
		return fmt.Errorf("failed to delete broker: %w", err)
	}

	return nil
}

// determineUserRole maps broker role to user role
func determineUserRole(brokerRole string) string {
	switch brokerRole {
	case "admin", "broker_admin":
		return "admin"
	case "manager":
		return "manager"
	default:
		return "admin" // Default to admin for safety
	}
}

func printStats(stats *MigrationStats) {
	log.Println("\n===========================================")
	log.Println("MIGRATION STATISTICS")
	log.Println("===========================================")
	log.Printf("Total Brokers Found:    %d\n", stats.TotalBrokers)
	log.Printf("Real Brokers (CRECI):   %d\n", stats.RealBrokers)
	log.Printf("Admin Users (No CRECI): %d\n", stats.AdminUsers)
	log.Printf("Migrated to /users:     %d\n", stats.Migrated)
	log.Printf("Kept in /brokers:       %d\n", stats.Skipped)
	log.Printf("Errors:                 %d\n", stats.Errors)
	log.Println("===========================================")

	if len(stats.ErrorDetails) > 0 {
		log.Println("\nERROR DETAILS:")
		for i, err := range stats.ErrorDetails {
			log.Printf("%d. %s\n", i+1, err)
		}
	}
}
