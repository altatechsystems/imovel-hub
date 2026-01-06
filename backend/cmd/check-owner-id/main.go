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
		serviceAccountPath = "config/firebase-adminsdk.json"
		if _, err := os.Stat(serviceAccountPath); os.IsNotExist(err) {
			log.Fatal("Firebase credentials not found")
		}
	}

	opt := option.WithCredentialsFile(serviceAccountPath)
	projectID := "ecosistema-imob-dev"
	databaseID := "imob-dev"

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, opt)
	if err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer client.Close()

	tenantID := "bd71c02b-5fa5-43df-8b46-e1df2206f1ef"
	ownerID := "e0887c7e-f347-4ccd-840c-e98badd97d1f"

	// First, let's check all possible locations
	fmt.Printf("🔍 Checking properties for owner: %s\n\n", ownerID)

	// First, just check if we can access the properties collection at all
	fmt.Println("🏠 Checking if properties collection exists...")
	allPropsQuery := client.Collection("properties").Limit(3)
	allPropsIter := allPropsQuery.Documents(ctx)
	defer allPropsIter.Stop()

	totalProps := 0
	for {
		doc, err := allPropsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("❌ Error accessing properties: %v\n", err)
			break
		}
		totalProps++
		data := doc.Data()
		fmt.Printf("   Sample property %d: %s (owner_id=%v)\n", totalProps, data["reference"], data["owner_id"])
	}

	if totalProps == 0 {
		fmt.Println("⚠️  No properties found at all!")
		return
	}

	fmt.Printf("\n✅ Found %d properties total\n\n", totalProps)

	// Now check properties with this specific owner_id
	fmt.Println("🔍 Checking properties for specific owner...")
	query := client.Collection("properties").
		Where("owner_id", "==", ownerID).
		Limit(10)

	iter := query.Documents(ctx)
	defer iter.Stop()

	propertiesFound := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("❌ Error: %v\n", err)
			break
		}

		propertiesFound++
		data := doc.Data()
		fmt.Printf("✅ Found property %d:\n", propertiesFound)
		fmt.Printf("   ID: %s\n", doc.Ref.ID)
		fmt.Printf("   Reference: %v\n", data["reference"])
		fmt.Printf("   Title: %v\n", data["title"])
		fmt.Printf("   Owner ID: %v\n\n", data["owner_id"])
	}

	if propertiesFound == 0 {
		fmt.Printf("\n⚠️  No properties found for owner: %s\n", ownerID)
		fmt.Println("\n🔍 Let's check what owners actually exist...")

		// Sample some properties to see what owner_ids exist
		sampleQuery := client.Collection("properties").Where("tenant_id", "==", tenantID).Limit(5)
		sampleIter := sampleQuery.Documents(ctx)
		defer sampleIter.Stop()

		fmt.Println("\n📋 Sample of properties in database:")
		sampleCount := 0
		for {
			doc, err := sampleIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				break
			}
			sampleCount++
			data := doc.Data()
			fmt.Printf("   Property %d: ref=%v, owner_id=%v\n", sampleCount, data["reference"], data["owner_id"])
		}
	} else {
		fmt.Printf("\n✅ SUCCESS: Found %d properties for this owner!\n", propertiesFound)
	}
}

func countDocs(ctx context.Context, client *firestore.Client, query firestore.Query) int {
	iter := query.Limit(10).Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("   ❌ Error: %v\n", err)
			return 0
		}
		count++

		// Show first document as sample
		if count == 1 {
			data := doc.Data()
			fmt.Printf("   📄 Sample doc ID: %s\n", doc.Ref.ID)
			fmt.Printf("      Reference: %v\n", data["reference"])
			fmt.Printf("      Owner ID: %v\n", data["owner_id"])
			fmt.Printf("      Name: %v\n", data["name"])
		}
	}
	fmt.Printf("   Found: %d documents\n\n", count)
	return count
}

func countDocsFromCollection(ctx context.Context, client *firestore.Client, col *firestore.CollectionRef) int {
	iter := col.Limit(10).Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("   ❌ Error: %v\n", err)
			return 0
		}
		count++

		// Show first document as sample
		if count == 1 {
			data := doc.Data()
			fmt.Printf("   📄 Sample doc ID: %s\n", doc.Ref.ID)
			fmt.Printf("      Reference: %v\n", data["reference"])
			fmt.Printf("      Owner ID: %v\n", data["owner_id"])
			fmt.Printf("      Name: %v\n", data["name"])
		}
	}
	fmt.Printf("   Found: %d documents\n\n", count)
	return count
}
