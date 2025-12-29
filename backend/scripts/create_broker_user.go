package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Initialize Firebase
	opt := option.WithCredentialsFile("../firebase-credentials.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}

	// Get Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Auth client: %v", err)
	}

	// Get Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// Tenant ID - ALTATECH Imóveis
	tenantID := "bd71c02b-5fa5-43df-8b46-a1df2206f1ef"
	email := "corretor@altatech.imob"
	password := "senha123"
	name := "Corretor ALTATECH"

	// Create Firebase user
	userParams := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(name).
		EmailVerified(true)

	userRecord, err := authClient.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating Firebase user: %v", err)
		log.Println("Trying to get existing user...")
		userRecord, err = authClient.GetUserByEmail(ctx, email)
		if err != nil {
			log.Fatalf("Error getting user: %v", err)
		}
		log.Printf("Found existing user: %s", userRecord.UID)
	} else {
		log.Printf("Created Firebase user: %s", userRecord.UID)
	}

	// Create broker in Firestore
	brokerID := uuid.New().String()
	now := time.Now()

	broker := models.Broker{
		ID:           brokerID,
		TenantID:     tenantID,
		FirebaseUID:  userRecord.UID,
		Name:         name,
		Email:        email,
		Phone:        "+55 11 99999-9999",
		Role:         models.BrokerRoleAdmin, // Admin role
		IsActive:     true,
		CanManageAll: true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = firestoreClient.Collection("brokers").Doc(brokerID).Set(ctx, broker)
	if err != nil {
		log.Fatalf("Error creating broker: %v", err)
	}

	fmt.Printf("\n✅ Broker created successfully!\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Tenant ID: %s\n", tenantID)
	fmt.Printf("Broker ID: %s\n", brokerID)
	fmt.Printf("Firebase UID: %s\n", userRecord.UID)
}
