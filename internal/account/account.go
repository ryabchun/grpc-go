package account

import (
	"database/sql"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
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

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func (s *Service) CreateAccount(username, email, password string) (*Account, error) {
	if !isValidEmail(email) {
		return nil, errors.New("invalid email")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	acc := &Account{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Status:   "ACTIVE",
	}
	id, err := s.repo.Create(acc)
	if err != nil {
		return nil, err
	}
	acc.ID = id
	return acc, nil
}

func (s *Service) GetAccount(id int64) (*Account, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateAccount(id int64, username, email, password, status string) (*Account, error) {
	acc, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	acc.Username = username
	acc.Email = email
	acc.Status = status
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		acc.Password = string(hashedPassword)
	}
	err = s.repo.Update(acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *Service) DeleteAccount(id int64) error {
	return s.repo.Delete(id)
}
