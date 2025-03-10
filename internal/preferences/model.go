package preferences

import "time"

type Preference struct {
	ID            int64
	AccountID     int64
	Theme         string
	Notifications bool
	Locale        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
