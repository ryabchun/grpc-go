package orders

import "time"

type OrderItem struct {
	Product  string
	Quantity int32
	Price    float32
}

type Order struct {
	ID        int64
	AccountID int64
	Items     []OrderItem
	Total     float32
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
