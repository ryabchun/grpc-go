package account_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"grpc-go/internal/account"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(acc *account.Account) (int64, error) {
	args := m.Called(acc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) Get(id int64) (*account.Account, error) {
	args := m.Called(id)
	if acc, ok := args.Get(0).(*account.Account); ok {
		return acc, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(acc *account.Account) error {
	args := m.Called(acc)
	return args.Error(0)
}

func (m *MockRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := account.NewServiceWithRepository(mockRepo)

	username := "testuser"
	email := "test@example.com"
	password := "password123"
	expectedID := int64(1)

	mockRepo.
		On("Create", mock.AnythingOfType("*account.Account")).
		Return(expectedID, nil)

	acc, err := svc.CreateAccount(username, email, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, acc.ID)
	assert.Equal(t, username, acc.Username)
	assert.Equal(t, email, acc.Email)
	assert.Equal(t, "ACTIVE", acc.Status)
	mockRepo.AssertExpectations(t)
}

func TestGetAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := account.NewServiceWithRepository(mockRepo)
	accountData := &account.Account{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("Get", int64(1)).Return(accountData, nil)

	acc, err := svc.GetAccount(1)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", acc.Username)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := account.NewServiceWithRepository(mockRepo)
	original := &account.Account{
		ID:       1,
		Username: "olduser",
		Email:    "old@example.com",
		Password: "oldhash",
		Status:   "ACTIVE",
	}
	mockRepo.On("Get", int64(1)).Return(original, nil)
	mockRepo.
		On("Update", mock.MatchedBy(func(acc *account.Account) bool {
			return acc.Username == "newuser" && acc.Email == "new@example.com" && acc.Status == "SUSPENDED"
		})).
		Return(nil)

	result, err := svc.UpdateAccount(1, "newuser", "new@example.com", "newpassword", "SUSPENDED")
	assert.NoError(t, err)
	assert.Equal(t, "newuser", result.Username)
	assert.Equal(t, "new@example.com", result.Email)
	assert.Equal(t, "SUSPENDED", result.Status)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := account.NewServiceWithRepository(mockRepo)
	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.DeleteAccount(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_InvalidEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := account.NewServiceWithRepository(mockRepo)

	_, err := svc.CreateAccount("user", "invalid-email", "password")
	assert.Error(t, err)
}
