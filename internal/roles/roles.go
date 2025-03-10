package roles

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

func (s *Service) CreateRole(name string, permissions []string) (*Role, error) {
	role := &Role{
		Name:        name,
		Permissions: permissions,
	}
	id, err := s.repo.Create(role)
	if err != nil {
		return nil, err
	}
	role.ID = id
	return role, nil
}

func (s *Service) GetRole(id int64) (*Role, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateRole(id int64, name string, permissions []string) (*Role, error) {
	role, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	role.Name = name
	role.Permissions = permissions
	err = s.repo.Update(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *Service) DeleteRole(id int64) error {
	return s.repo.Delete(id)
}
