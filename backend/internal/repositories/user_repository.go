package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserRepository handles database operations for administrative users
type UserRepository struct {
	client *firestore.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *firestore.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	if user.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := r.getUsersCollection(user.TenantID)
	_, err := collection.Doc(user.ID).Set(ctx, user)
	return err
}

// Get retrieves a user by ID
func (r *UserRepository) Get(ctx context.Context, tenantID, userID string) (*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	collection := r.getUsersCollection(tenantID)
	doc, err := collection.Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user: %w", err)
	}

	user.ID = doc.Ref.ID
	return &user, nil
}

// GetByEmail retrieves a user by email within a tenant
func (r *UserRepository) GetByEmail(ctx context.Context, tenantID, email string) (*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	collection := r.getUsersCollection(tenantID)
	query := collection.Where("email", "==", email).Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user: %w", err)
	}

	user.ID = doc.Ref.ID
	return &user, nil
}

// GetByFirebaseUID retrieves a user by Firebase UID within a tenant
func (r *UserRepository) GetByFirebaseUID(ctx context.Context, tenantID, firebaseUID string) (*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if firebaseUID == "" {
		return nil, fmt.Errorf("firebase UID is required")
	}

	collection := r.getUsersCollection(tenantID)
	query := collection.Where("firebase_uid", "==", firebaseUID).Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user by firebase UID: %w", err)
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user: %w", err)
	}

	user.ID = doc.Ref.ID
	return &user, nil
}

// List retrieves all users for a tenant
func (r *UserRepository) List(ctx context.Context, tenantID string) ([]*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}

	collection := r.getUsersCollection(tenantID)
	iter := collection.Documents(ctx)
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %w", err)
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to parse user: %w", err)
		}

		user.ID = doc.Ref.ID
		users = append(users, &user)
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, tenantID, userID string, updates map[string]interface{}) error {
	if tenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	updates["updated_at"] = time.Now()

	collection := r.getUsersCollection(tenantID)
	_, err := collection.Doc(userID).Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, tenantID, userID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	collection := r.getUsersCollection(tenantID)
	_, err := collection.Doc(userID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListByRole retrieves users by role
func (r *UserRepository) ListByRole(ctx context.Context, tenantID, role string) ([]*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}

	collection := r.getUsersCollection(tenantID)
	query := collection.Where("role", "==", role)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %w", err)
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to parse user: %w", err)
		}

		user.ID = doc.Ref.ID
		users = append(users, &user)
	}

	return users, nil
}

// ListActive retrieves all active users for a tenant
func (r *UserRepository) ListActive(ctx context.Context, tenantID string) ([]*models.User, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}

	collection := r.getUsersCollection(tenantID)
	query := collection.Where("is_active", "==", true)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate users: %w", err)
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to parse user: %w", err)
		}

		user.ID = doc.Ref.ID
		users = append(users, &user)
	}

	return users, nil
}

// getUsersCollection returns the users collection for a tenant
func (r *UserRepository) getUsersCollection(tenantID string) *firestore.CollectionRef {
	return r.client.Collection("tenants").Doc(tenantID).Collection("users")
}
