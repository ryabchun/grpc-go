package orders

import (
	"database/sql"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{repo: NewSQLRepository(db)}
}

func NewServiceWithRepository(repo Repository) *Service {
	return &Service{repo: repo}
}

func calculateTotal(items []OrderItem) float32 {
	var total float32
	for _, item := range items {
		total += float32(item.Quantity) * item.Price
	}
	return total
}

func (s *Service) CreateOrder(accountID int64, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}
	total := calculateTotal(items)
	order := &Order{
		AccountID: accountID,
		Items:     items,
		Total:     total,
		Status:    "PENDING",
	}
	id, err := s.repo.Create(order)
	if err != nil {
		return nil, err
	}
	order.ID = id
	return order, nil
}

func (s *Service) GetOrder(id int64) (*Order, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateOrder(id int64, items []OrderItem, status string) (*Order, error) {
	order, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}
	order.Items = items
	order.Total = calculateTotal(items)
	order.Status = status
	err = s.repo.Update(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) DeleteOrder(id int64) error {
	return s.repo.Delete(id)
}
