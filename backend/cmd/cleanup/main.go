package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Initialize Firebase
	serviceAccountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountPath == "" {
		// Try default path
		serviceAccountPath = "config/firebase-adminsdk.json"
		if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
			log.Fatal("Firebase credentials not found. Set GOOGLE_APPLICATION_CREDENTIALS or place config/firebase-adminsdk.json")
		}
	}

	// Initialize Firestore client directly with named database (same as server)
	opt := option.WithCredentialsFile(serviceAccountPath)
	projectID := "ecosistema-imob-dev"
	databaseID := "imob-dev"

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, opt)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
	}
	defer client.Close()

	tenantID := "bd71c02b-5fa5-43df-8b46-e1df2206f1ef"

	fmt.Println("🗑️  Starting database cleanup...")
	fmt.Printf("   Tenant ID: %s\n\n", tenantID)

	// Collections to clean
	collections := []string{
		"properties",
		"owners", // Old root collection
		"listings",
		"import_batches",
		"import_errors",
		"property_broker_roles",
	}

	// Clean root collections
	for _, collectionName := range collections {
		fmt.Printf("🧹 Cleaning collection: %s\n", collectionName)
		if err := deleteCollection(ctx, client, collectionName, tenantID); err != nil {
			log.Printf("   ⚠️  Error cleaning %s: %v", collectionName, err)
		} else {
			fmt.Printf("   ✅ Cleaned %s\n", collectionName)
		}
	}

	// Clean tenant-scoped collections
	fmt.Printf("\n🧹 Cleaning tenant-scoped collections\n")

	// First, check if tenant document exists
	tenantDoc, err := client.Collection("tenants").Doc(tenantID).Get(ctx)
	if err != nil {
		fmt.Printf("   ⚠️  Tenant document doesn't exist or error: %v\n", err)
		fmt.Printf("   ℹ️  Trying direct collection access instead...\n")

		// Try direct collection access
		tenantCollections := []string{"owners", "properties", "listings"}
		for _, collectionName := range tenantCollections {
			fullPath := fmt.Sprintf("tenants/%s/%s", tenantID, collectionName)
			fmt.Printf("🧹 Cleaning: %s (direct access)\n", fullPath)
			if err := deleteCollectionDirect(ctx, client, fullPath); err != nil {
				log.Printf("   ⚠️  Error cleaning %s: %v", fullPath, err)
			} else {
				fmt.Printf("   ✅ Cleaned %s\n", fullPath)
			}
		}
	} else {
		fmt.Printf("   ✅ Tenant document exists: %s\n", tenantDoc.Ref.ID)

		// Access tenant document and its subcollections
		tenantCollections := []string{"owners", "properties", "listings"}
		for _, collectionName := range tenantCollections {
			fullPath := fmt.Sprintf("tenants/%s/%s", tenantID, collectionName)
			fmt.Printf("🧹 Cleaning: %s\n", fullPath)
			if err := deleteTenantSubcollection(ctx, client, tenantID, collectionName); err != nil {
				log.Printf("   ⚠️  Error cleaning %s: %v", fullPath, err)
			} else {
				fmt.Printf("   ✅ Cleaned %s\n", fullPath)
			}
		}
	}

	fmt.Println("\n✅ Database cleanup completed!")
	fmt.Println("   You can now do a fresh import.")
}

// deleteCollection deletes all documents in a collection for a specific tenant
func deleteCollection(ctx context.Context, client *firestore.Client, collectionName, tenantID string) error {
	col := client.Collection(collectionName)

	// Get ALL documents (don't filter by tenant_id)
	iter := col.Documents(ctx)
	defer iter.Stop()

	batch := client.Batch()
	count := 0
	batchSize := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating documents: %w", err)
		}

		batch.Delete(doc.Ref)
		batchSize++
		count++

		// Commit batch every 500 documents
		if batchSize >= 500 {
			if _, err := batch.Commit(ctx); err != nil {
				return fmt.Errorf("error committing batch: %w", err)
			}
			batch = client.Batch()
			batchSize = 0
			fmt.Printf("   📦 Deleted %d documents...\n", count)
		}
	}

	// Commit remaining documents
	if batchSize > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			return fmt.Errorf("error committing final batch: %w", err)
		}
	}

	fmt.Printf("   🗑️  Total deleted: %d documents\n", count)
	return nil
}

// deleteCollectionDirect deletes all documents using direct collection path
func deleteCollectionDirect(ctx context.Context, client *firestore.Client, collectionPath string) error {
	col := client.Collection(collectionPath)

	iter := col.Documents(ctx)
	defer iter.Stop()

	batch := client.Batch()
	count := 0
	batchSize := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating documents: %w", err)
		}

		batch.Delete(doc.Ref)
		batchSize++
		count++

		// Commit batch every 500 documents
		if batchSize >= 500 {
			if _, err := batch.Commit(ctx); err != nil {
				return fmt.Errorf("error committing batch: %w", err)
			}
			batch = client.Batch()
			batchSize = 0
			fmt.Printf("   📦 Deleted %d documents...\n", count)
		}
	}

	// Commit remaining documents
	if batchSize > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			return fmt.Errorf("error committing final batch: %w", err)
		}
	}

	fmt.Printf("   🗑️  Total deleted: %d documents\n", count)
	return nil
}

// deleteTenantSubcollection deletes all documents in a tenant subcollection
func deleteTenantSubcollection(ctx context.Context, client *firestore.Client, tenantID, collectionName string) error {
	// Access subcollection through tenant document
	col := client.Collection("tenants").Doc(tenantID).Collection(collectionName)
	fmt.Printf("   🔍 Accessing: tenants/%s/%s\n", tenantID, collectionName)

	iter := col.Documents(ctx)
	defer iter.Stop()

	batch := client.Batch()
	count := 0
	batchSize := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating documents: %w", err)
		}

		batch.Delete(doc.Ref)
		batchSize++
		count++

		// Commit batch every 500 documents
		if batchSize >= 500 {
			if _, err := batch.Commit(ctx); err != nil {
				return fmt.Errorf("error committing batch: %w", err)
			}
			batch = client.Batch()
			batchSize = 0
			fmt.Printf("   📦 Deleted %d documents...\n", count)
		}
	}

	// Commit remaining documents
	if batchSize > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			return fmt.Errorf("error committing final batch: %w", err)
		}
	}

	fmt.Printf("   🗑️  Total deleted: %d documents\n", count)
	return nil
}

// deleteSubcollection deletes all documents in a subcollection
func deleteSubcollection(ctx context.Context, client *firestore.Client, parentPath, collectionName string) error {
	fullPath := fmt.Sprintf("%s/%s", parentPath, collectionName)
	fmt.Printf("   🔍 Accessing collection path: %s\n", fullPath)

	col := client.Collection(fullPath)

	iter := col.Documents(ctx)
	defer iter.Stop()

	batch := client.Batch()
	count := 0
	batchSize := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating documents: %w", err)
		}

		batch.Delete(doc.Ref)
		batchSize++
		count++

		// Commit batch every 500 documents
		if batchSize >= 500 {
			if _, err := batch.Commit(ctx); err != nil {
				return fmt.Errorf("error committing batch: %w", err)
			}
			batch = client.Batch()
			batchSize = 0
			fmt.Printf("   📦 Deleted %d documents...\n", count)
		}
	}

	// Commit remaining documents
	if batchSize > 0 {
		if _, err := batch.Commit(ctx); err != nil {
			return fmt.Errorf("error committing final batch: %w", err)
		}
	}

	fmt.Printf("   🗑️  Total deleted: %d documents\n", count)
	return nil
}
