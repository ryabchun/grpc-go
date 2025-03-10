package roles

import "time"

type Role struct {
	ID          int64
	Name        string
	Permissions []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
