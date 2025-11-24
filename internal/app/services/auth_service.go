package services

import (
	"backend/internal/domain/models"
)

type AuthService interface {
	Register(user *models.User) (*models.User, error)
	Login(username, password string) (string, *models.User, error)
}