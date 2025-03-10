package account

import "time"

type Account struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
