package preferences_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"grpc-go/internal/preferences"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(pref *preferences.Preference) (int64, error) {
	args := m.Called(pref)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) Get(id int64) (*preferences.Preference, error) {
	args := m.Called(id)
	if pref, ok := args.Get(0).(*preferences.Preference); ok {
		return pref, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(pref *preferences.Preference) error {
	args := m.Called(pref)
	return args.Error(0)
}

func (m *MockRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreatePreference(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := preferences.NewServiceWithRepository(mockRepo)
	mockRepo.
		On("Create", mock.AnythingOfType("*preferences.Preference")).
		Return(int64(1), nil)

	pref, err := svc.CreatePreference(1, "dark", true, "en-US")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), pref.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetPreference(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := preferences.NewServiceWithRepository(mockRepo)
	prefData := &preferences.Preference{
		ID:            1,
		AccountID:     1,
		Theme:         "light",
		Notifications: false,
		Locale:        "en-GB",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	mockRepo.On("Get", int64(1)).Return(prefData, nil)

	pref, err := svc.GetPreference(1)
	assert.NoError(t, err)
	assert.Equal(t, "en-GB", pref.Locale)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePreference(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := preferences.NewServiceWithRepository(mockRepo)
	original := &preferences.Preference{
		ID:            1,
		AccountID:     1,
		Theme:         "dark",
		Notifications: true,
		Locale:        "en-US",
	}
	mockRepo.On("Get", int64(1)).Return(original, nil)
	mockRepo.
		On("Update", mock.MatchedBy(func(p *preferences.Preference) bool {
			return p.Theme == "light" && p.Locale == "en-CA"
		})).
		Return(nil)

	updated, err := svc.UpdatePreference(1, "light", false, "en-CA")
	assert.NoError(t, err)
	assert.Equal(t, "light", updated.Theme)
	assert.Equal(t, "en-CA", updated.Locale)
	mockRepo.AssertExpectations(t)
}

func TestDeletePreference(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := preferences.NewServiceWithRepository(mockRepo)
	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.DeletePreference(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
