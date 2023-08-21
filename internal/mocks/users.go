package mocks

import (
	"time"

	"github.com/ahmadyogi543/snippetbox/internal/models"
)

type UserModel struct{}

func (um *UserModel) Get(id int) (*models.User, error) {
	if id == 1 {
		return &models.User{
			ID:      1,
			Name:    "Ahmad Yogi",
			Email:   "ayogi@snippetbox.sh",
			Created: time.Now(),
		}, nil
	}

	return nil, models.ErrNoRecord
}

func (um *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
	return nil
}

func (um *UserModel) Insert(name, email, password string) error {
	switch email {
	case "duplicate@snippetbox.sh":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	if email == "ayogi@snippetbox.sh" && password == "12345678" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (um *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
