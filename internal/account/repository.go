package account

import (
	"database/sql"
)

type Repository interface {
	Create(account *Account) (int64, error)
	Get(id int64) (*Account, error)
	Update(account *Account) error
	Delete(id int64) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(account *Account) (int64, error) {
	query := `
		INSERT INTO accounts (username, email, password, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id
	`
	var id int64
	err := r.db.QueryRow(query, account.Username, account.Email, account.Password, account.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SQLRepository) Get(id int64) (*Account, error) {
	query := `SELECT id, username, email, status FROM accounts WHERE id = $1`
	var account Account
	err := r.db.QueryRow(query, id).Scan(&account.ID, &account.Username, &account.Email, &account.Status)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *SQLRepository) Update(account *Account) error {
	query := `UPDATE accounts SET username = $1, email = $2, password = $3, status = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.db.Exec(query, account.Username, account.Email, account.Password, account.Status, account.ID)
	return err
}

func (r *SQLRepository) Delete(id int64) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
