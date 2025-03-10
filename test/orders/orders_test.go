package orders_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"grpc-go/internal/orders"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(order *orders.Order) (int64, error) {
	args := m.Called(order)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) Get(id int64) (*orders.Order, error) {
	args := m.Called(id)
	if order, ok := args.Get(0).(*orders.Order); ok {
		return order, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(order *orders.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	// Assume we added NewServiceWithRepository to orders service.
	svc := orders.NewServiceWithRepository(mockRepo)
	items := []orders.OrderItem{
		{Product: "Product1", Quantity: 2, Price: 10.0},
		{Product: "Product2", Quantity: 1, Price: 20.0},
	}
	expectedTotal := float32(2*10.0 + 1*20.0)
	mockRepo.
		On("Create", mock.MatchedBy(func(o *orders.Order) bool {
			return o.AccountID == 1 && o.Total == expectedTotal && o.Status == "PENDING"
		})).
		Return(int64(1), nil)

	order, err := svc.CreateOrder(1, items)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), order.ID)
	assert.Equal(t, expectedTotal, order.Total)
	mockRepo.AssertExpectations(t)
}

func TestGetOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := orders.NewServiceWithRepository(mockRepo)
	orderData := &orders.Order{
		ID:        1,
		AccountID: 1,
		Total:     30.0,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("Get", int64(1)).Return(orderData, nil)

	order, err := svc.GetOrder(1)
	assert.NoError(t, err)
	assert.Equal(t, orderData.Status, order.Status)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := orders.NewServiceWithRepository(mockRepo)
	originalOrder := &orders.Order{
		ID:        1,
		AccountID: 1,
		Total:     30.0,
		Status:    "PENDING",
	}
	mockRepo.On("Get", int64(1)).Return(originalOrder, nil)
	newItems := []orders.OrderItem{
		{Product: "Product1", Quantity: 3, Price: 10.0},
	}
	expectedTotal := float32(3 * 10.0)
	mockRepo.
		On("Update", mock.MatchedBy(func(o *orders.Order) bool {
			return o.Total == expectedTotal && o.Status == "COMPLETED"
		})).
		Return(nil)

	updated, err := svc.UpdateOrder(1, newItems, "COMPLETED")
	assert.NoError(t, err)
	assert.Equal(t, "COMPLETED", updated.Status)
	mockRepo.AssertExpectations(t)
}

func TestDeleteOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := orders.NewServiceWithRepository(mockRepo)
	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.DeleteOrder(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
