package preferences

import (
	"database/sql"
)

type Repository interface {
	Create(pref *Preference) (int64, error)
	Get(id int64) (*Preference, error)
	Update(pref *Preference) error
	Delete(id int64) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(pref *Preference) (int64, error) {
	query := `
		INSERT INTO preferences (account_id, theme, notifications, locale, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id
	`
	var id int64
	err := r.db.QueryRow(query, pref.AccountID, pref.Theme, pref.Notifications, pref.Locale).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SQLRepository) Get(id int64) (*Preference, error) {
	query := `SELECT id, account_id, theme, notifications, locale FROM preferences WHERE id = $1`
	var pref Preference
	err := r.db.QueryRow(query, id).Scan(&pref.ID, &pref.AccountID, &pref.Theme, &pref.Notifications, &pref.Locale)
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *SQLRepository) Update(pref *Preference) error {
	query := `UPDATE preferences SET theme = $1, notifications = $2, locale = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(query, pref.Theme, pref.Notifications, pref.Locale, pref.ID)
	return err
}

func (r *SQLRepository) Delete(id int64) error {
	query := `DELETE FROM preferences WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
