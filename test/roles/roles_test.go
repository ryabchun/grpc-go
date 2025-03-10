package roles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"grpc-go/internal/roles"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(role *roles.Role) (int64, error) {
	args := m.Called(role)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) Get(id int64) (*roles.Role, error) {
	args := m.Called(id)
	if role, ok := args.Get(0).(*roles.Role); ok {
		return role, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(role *roles.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateRole(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := roles.NewServiceWithRepository(mockRepo)
	roleName := "admin"
	permissions := []string{"orders:write", "accounts:manage"}
	expectedID := int64(1)
	mockRepo.
		On("Create", mock.AnythingOfType("*roles.Role")).
		Return(expectedID, nil)

	role, err := svc.CreateRole(roleName, permissions)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, role.ID)
	assert.Equal(t, roleName, role.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetRole(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := roles.NewServiceWithRepository(mockRepo)
	roleData := &roles.Role{
		ID:          1,
		Name:        "user",
		Permissions: []string{"orders:read"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo.On("Get", int64(1)).Return(roleData, nil)

	role, err := svc.GetRole(1)
	assert.NoError(t, err)
	assert.Equal(t, "user", role.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRole(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := roles.NewServiceWithRepository(mockRepo)
	original := &roles.Role{
		ID:          1,
		Name:        "user",
		Permissions: []string{"orders:read"},
	}
	mockRepo.On("Get", int64(1)).Return(original, nil)
	mockRepo.
		On("Update", mock.MatchedBy(func(r *roles.Role) bool {
			return r.Name == "superuser" && len(r.Permissions) == 2
		})).
		Return(nil)

	updated, err := svc.UpdateRole(1, "superuser", []string{"orders:read", "orders:write"})
	assert.NoError(t, err)
	assert.Equal(t, "superuser", updated.Name)
	mockRepo.AssertExpectations(t)
}

func TestDeleteRole(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := roles.NewServiceWithRepository(mockRepo)
	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.DeleteRole(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
