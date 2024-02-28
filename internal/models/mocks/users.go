package mocks

import "github.com/svidlak/lets-go/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (*models.User, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return &models.User{}, nil
	}

	return &models.User{}, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
