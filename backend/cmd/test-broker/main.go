package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	projectID := "ecosistema-imob-dev"
	databaseID := "imob-dev"
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if credentialsPath == "" {
		credentialsPath = "./config/firebase-adminsdk.json"
	}

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	tenantID := "bd71c02b-5fa5-43df-8b46-a1df2206f1ef"
	brokerID := "CifbkWdJYelWvbijHI7Q"

	fmt.Printf("🔍 Buscando broker %s do tenant %s\n", brokerID, tenantID)

	// Get broker document
	docPath := fmt.Sprintf("tenants/%s/brokers/%s", tenantID, brokerID)
	doc, err := client.Doc(docPath).Get(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to get broker document: %v", err)
	}

	fmt.Println("✅ Document exists")
	fmt.Println("📄 Raw data:")
	fmt.Println(doc.Data())

	// Try to unmarshal to struct
	var broker models.Broker
	if err := doc.DataTo(&broker); err != nil {
		log.Fatalf("❌ Failed to unmarshal broker: %v", err)
	}

	broker.ID = doc.Ref.ID
	fmt.Println("\n✅ Successfully unmarshaled to struct:")
	fmt.Printf("   ID: %s\n", broker.ID)
	fmt.Printf("   Name: %s\n", broker.Name)
	fmt.Printf("   Email: %s\n", broker.Email)
	fmt.Printf("   CRECI: %s\n", broker.CRECI)
	fmt.Printf("   Role: %s\n", broker.Role)
	fmt.Printf("   IsActive: %v\n", broker.IsActive)
	fmt.Printf("   FirebaseUID: %s\n", broker.FirebaseUID)
}
