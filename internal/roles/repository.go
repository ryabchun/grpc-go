package roles

import (
	"database/sql"
	"strings"
)

type Repository interface {
	Create(role *Role) (int64, error)
	Get(id int64) (*Role, error)
	Update(role *Role) error
	Delete(id int64) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(role *Role) (int64, error) {
	perms := strings.Join(role.Permissions, ",")
	query := `
		INSERT INTO roles (name, permissions, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING id
	`
	var id int64
	err := r.db.QueryRow(query, role.Name, perms).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SQLRepository) Get(id int64) (*Role, error) {
	query := `SELECT id, name, permissions FROM roles WHERE id = $1`
	var role Role
	var perms string
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &perms)
	if err != nil {
		return nil, err
	}
	role.Permissions = strings.Split(perms, ",")
	return &role, nil
}

func (r *SQLRepository) Update(role *Role) error {
	perms := strings.Join(role.Permissions, ",")
	query := `UPDATE roles SET name = $1, permissions = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, role.Name, perms, role.ID)
	return err
}

func (r *SQLRepository) Delete(id int64) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
