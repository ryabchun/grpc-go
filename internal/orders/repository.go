package orders

import (
	"database/sql"
)

type Repository interface {
	Create(order *Order) (int64, error)
	Get(id int64) (*Order, error)
	Update(order *Order) error
	Delete(id int64) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(order *Order) (int64, error) {
	query := `
		INSERT INTO orders (account_id, total, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id
	`
	var id int64
	err := r.db.QueryRow(query, order.AccountID, order.Total, order.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO order_items (order_id, product, quantity, price)
			VALUES ($1, $2, $3, $4)
		`
		_, err := r.db.Exec(itemQuery, id, item.Product, item.Quantity, item.Price)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (r *SQLRepository) Get(id int64) (*Order, error) {
	query := `SELECT id, account_id, total, status FROM orders WHERE id = $1`
	var order Order
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.AccountID, &order.Total, &order.Status)
	if err != nil {
		return nil, err
	}

	itemsQuery := `SELECT product, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.Product, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return &order, nil
}

func (r *SQLRepository) Update(order *Order) error {
	query := `UPDATE orders SET total = $1, status = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, order.Total, order.Status, order.ID)
	if err != nil {
		return err
	}

	delQuery := `DELETE FROM order_items WHERE order_id = $1`
	_, err = r.db.Exec(delQuery, order.ID)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO order_items (order_id, product, quantity, price)
			VALUES ($1, $2, $3, $4)
		`
		_, err := r.db.Exec(itemQuery, order.ID, item.Product, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM order_items WHERE order_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM orders WHERE id = $1`, id)
	return err
}
