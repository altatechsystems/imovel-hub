package services

import (
	"context"
	"testing"
	"time"

	"github.com/altatech/ecosistema-imob/backend/internal/models"
	"github.com/altatech/ecosistema-imob/backend/internal/repositories"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users       map[string]*models.User
	emailIndex  map[string]string // email -> userID
	uidIndex    map[string]string // firebaseUID -> userID
	createError error
	getError    error
	updateError error
	deleteError error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:      make(map[string]*models.User),
		emailIndex: make(map[string]string),
		uidIndex:   make(map[string]string),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.createError != nil {
		return m.createError
	}
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user.ID
	m.uidIndex[user.FirebaseUID] = user.ID
	return nil
}

func (m *MockUserRepository) Get(ctx context.Context, tenantID, userID string) (*models.User, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	user, ok := m.users[userID]
	if !ok {
		return nil, repositories.ErrNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, tenantID, email string) (*models.User, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	userID, ok := m.emailIndex[email]
	if !ok {
		return nil, repositories.ErrNotFound
	}
	return m.users[userID], nil
}

func (m *MockUserRepository) GetByFirebaseUID(ctx context.Context, tenantID, firebaseUID string) (*models.User, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	userID, ok := m.uidIndex[firebaseUID]
	if !ok {
		return nil, repositories.ErrNotFound
	}
	return m.users[userID], nil
}

func (m *MockUserRepository) List(ctx context.Context, tenantID string) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		if user.TenantID == tenantID {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) ListByRole(ctx context.Context, tenantID, role string) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		if user.TenantID == tenantID && user.Role == role {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) ListActive(ctx context.Context, tenantID string) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		if user.TenantID == tenantID && user.IsActive {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, tenantID, userID string, updates map[string]interface{}) error {
	if m.updateError != nil {
		return m.updateError
	}
	user, ok := m.users[userID]
	if !ok {
		return repositories.ErrNotFound
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if email, ok := updates["email"].(string); ok {
		delete(m.emailIndex, user.Email)
		user.Email = email
		m.emailIndex[email] = userID
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if role, ok := updates["role"].(string); ok {
		user.Role = role
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		user.IsActive = isActive
	}
	if permissions, ok := updates["permissions"].([]string); ok {
		user.Permissions = permissions
	}

	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, tenantID, userID string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	user, ok := m.users[userID]
	if !ok {
		return repositories.ErrNotFound
	}
	delete(m.users, userID)
	delete(m.emailIndex, user.Email)
	delete(m.uidIndex, user.FirebaseUID)
	return nil
}

// MockTenantRepository is a mock implementation of TenantRepository for testing
type MockTenantRepository struct {
	tenants map[string]*models.Tenant
}

func NewMockTenantRepository() *MockTenantRepository {
	return &MockTenantRepository{
		tenants: make(map[string]*models.Tenant),
	}
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *MockTenantRepository) Get(ctx context.Context, tenantID string) (*models.Tenant, error) {
	tenant, ok := m.tenants[tenantID]
	if !ok {
		return nil, repositories.ErrNotFound
	}
	return tenant, nil
}

func (m *MockTenantRepository) Update(ctx context.Context, tenantID string, updates map[string]interface{}) error {
	return nil
}

func (m *MockTenantRepository) Delete(ctx context.Context, tenantID string) error {
	delete(m.tenants, tenantID)
	return nil
}

func (m *MockTenantRepository) List(ctx context.Context) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	for _, t := range m.tenants {
		tenants = append(tenants, t)
	}
	return tenants, nil
}

func (m *MockTenantRepository) ListActive(ctx context.Context) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	for _, t := range m.tenants {
		if t.IsActive {
			tenants = append(tenants, t)
		}
	}
	return tenants, nil
}

func (m *MockTenantRepository) GetBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	for _, t := range m.tenants {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, repositories.ErrNotFound
}

// MockActivityLogRepository is a mock implementation of ActivityLogRepository for testing
type MockActivityLogRepository struct {
	logs []*models.ActivityLog
}

func NewMockActivityLogRepository() *MockActivityLogRepository {
	return &MockActivityLogRepository{
		logs: make([]*models.ActivityLog, 0),
	}
}

func (m *MockActivityLogRepository) Create(ctx context.Context, log *models.ActivityLog) error {
	log.ID = "log-" + time.Now().Format("20060102150405")
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockActivityLogRepository) List(ctx context.Context, tenantID string, limit int) ([]*models.ActivityLog, error) {
	var filtered []*models.ActivityLog
	for _, log := range m.logs {
		if log.TenantID == tenantID {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

// Test CreateUser - Success
func TestCreateUser_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	// Create tenant first
	tenant := &models.Tenant{
		ID:       "tenant-1",
		Name:     "Test Tenant",
		IsActive: true,
	}
	mockTenantRepo.Create(context.Background(), tenant)

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Phone:       "11987654321",
		Role:        "admin",
	}

	err := service.CreateUser(context.Background(), user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was created
	created, err := mockUserRepo.Get(context.Background(), "tenant-1", "user-1")
	if err != nil {
		t.Errorf("Expected user to be created, got error %v", err)
	}
	if created.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", created.Email)
	}
}

// Test CreateUser - Missing Required Fields
func TestCreateUser_MissingTenantID(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	user := &models.User{
		ID:          "user-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
	}

	err := service.CreateUser(context.Background(), user)
	if err == nil {
		t.Error("Expected error for missing tenant_id, got nil")
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Email:       "john@example.com",
	}

	err := service.CreateUser(context.Background(), user)
	if err == nil {
		t.Error("Expected error for missing name, got nil")
	}
}

func TestCreateUser_MissingEmail(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
	}

	err := service.CreateUser(context.Background(), user)
	if err == nil {
		t.Error("Expected error for missing email, got nil")
	}
}

// Test CreateUser - Duplicate Email
func TestCreateUser_DuplicateEmail(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	// Create tenant first
	tenant := &models.Tenant{
		ID:       "tenant-1",
		Name:     "Test Tenant",
		IsActive: true,
	}
	mockTenantRepo.Create(context.Background(), tenant)

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	// Create first user
	user1 := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Role:        "admin",
	}
	service.CreateUser(context.Background(), user1)

	// Try to create second user with same email
	user2 := &models.User{
		ID:          "user-2",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-2",
		Name:        "Jane Admin",
		Email:       "john@example.com",
		Role:        "admin",
	}

	err := service.CreateUser(context.Background(), user2)
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}

// Test CreateUser - Invalid Role
func TestCreateUser_InvalidRole(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	// Create tenant first
	tenant := &models.Tenant{
		ID:       "tenant-1",
		Name:     "Test Tenant",
		IsActive: true,
	}
	mockTenantRepo.Create(context.Background(), tenant)

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Role:        "broker", // Invalid role for admin users
	}

	err := service.CreateUser(context.Background(), user)
	if err == nil {
		t.Error("Expected error for invalid role 'broker', got nil")
	}
}

// Test UpdateUser - Success
func TestUpdateUser_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	// Create user first
	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Role:        "admin",
	}
	mockUserRepo.Create(context.Background(), user)

	// Update user
	updates := map[string]interface{}{
		"name": "John Updated",
	}

	err := service.UpdateUser(context.Background(), "tenant-1", "user-1", updates)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	updated, _ := mockUserRepo.Get(context.Background(), "tenant-1", "user-1")
	if updated.Name != "John Updated" {
		t.Errorf("Expected name 'John Updated', got '%s'", updated.Name)
	}
}

// Test GrantPermission
func TestGrantPermission(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	// Create user
	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Role:        "manager",
		Permissions: []string{},
	}
	mockUserRepo.Create(context.Background(), user)

	// Grant permission
	err := service.GrantPermission(context.Background(), "tenant-1", "user-1", "properties.edit")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify permission was added
	updated, _ := mockUserRepo.Get(context.Background(), "tenant-1", "user-1")
	if len(updated.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(updated.Permissions))
	}
	if updated.Permissions[0] != "properties.edit" {
		t.Errorf("Expected permission 'properties.edit', got '%s'", updated.Permissions[0])
	}
}

// Test RevokePermission
func TestRevokePermission(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	// Create user with permission
	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Role:        "manager",
		Permissions: []string{"properties.edit"},
	}
	mockUserRepo.Create(context.Background(), user)

	// Revoke permission
	err := service.RevokePermission(context.Background(), "tenant-1", "user-1", "properties.edit")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify permission was removed
	updated, _ := mockUserRepo.Get(context.Background(), "tenant-1", "user-1")
	if len(updated.Permissions) != 0 {
		t.Errorf("Expected 0 permissions, got %d", len(updated.Permissions))
	}
}

// Test DeleteUser
func TestDeleteUser(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockTenantRepo := NewMockTenantRepository()
	mockActivityLogRepo := NewMockActivityLogRepository()

	service := NewUserService(mockUserRepo, mockTenantRepo, mockActivityLogRepo)

	// Create user
	user := &models.User{
		ID:          "user-1",
		TenantID:    "tenant-1",
		FirebaseUID: "firebase-uid-1",
		Name:        "John Admin",
		Email:       "john@example.com",
		Role:        "admin",
	}
	mockUserRepo.Create(context.Background(), user)

	// Delete user
	err := service.DeleteUser(context.Background(), "tenant-1", "user-1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was deleted
	_, err = mockUserRepo.Get(context.Background(), "tenant-1", "user-1")
	if err != repositories.ErrNotFound {
		t.Error("Expected user to be deleted")
	}
}
