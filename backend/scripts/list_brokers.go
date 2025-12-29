package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Initialize Firebase
	opt := option.WithCredentialsFile("../firebase-credentials.json")
	config := &firebase.Config{
		ProjectID:   "ecosistema-imob-dev",
		DatabaseURL: "https://ecosistema-imob-dev.firebaseio.com",
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}

	// Get Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// List all brokers
	fmt.Println("📋 Listing all brokers:")
	fmt.Println("---")

	iter := firestoreClient.Collection("brokers").Documents(ctx)
	count := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating brokers: %v", err)
		}

		data := doc.Data()
		count++
		fmt.Printf("\n%d. Broker ID: %s\n", count, doc.Ref.ID)
		fmt.Printf("   Name: %v\n", data["name"])
		fmt.Printf("   Email: %v\n", data["email"])
		fmt.Printf("   Tenant ID: %v\n", data["tenant_id"])
		fmt.Printf("   Firebase UID: %v\n", data["firebase_uid"])
		fmt.Printf("   Role: %v\n", data["role"])
		fmt.Printf("   Active: %v\n", data["is_active"])
	}

	fmt.Printf("\n---\nTotal brokers: %d\n", count)
}
