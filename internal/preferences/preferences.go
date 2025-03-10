package preferences

import (
	"database/sql"
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

func (s *Service) CreatePreference(accountID int64, theme string, notifications bool, locale string) (*Preference, error) {
	pref := &Preference{
		AccountID:     accountID,
		Theme:         theme,
		Notifications: notifications,
		Locale:        locale,
	}
	id, err := s.repo.Create(pref)
	if err != nil {
		return nil, err
	}
	pref.ID = id
	return pref, nil
}

func (s *Service) GetPreference(id int64) (*Preference, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdatePreference(id int64, theme string, notifications bool, locale string) (*Preference, error) {
	pref, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	pref.Theme = theme
	pref.Notifications = notifications
	pref.Locale = locale
	err = s.repo.Update(pref)
	if err != nil {
		return nil, err
	}
	return pref, nil
}

func (s *Service) DeletePreference(id int64) error {
	return s.repo.Delete(id)
}
