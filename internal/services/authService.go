package services

import (
	"github.com/ruziba3vich/exam/internal/models"
	"github.com/ruziba3vich/exam/internal/repositories"
)

type AuthService struct {
	storage repositories.IAuthenticationRepository
}

func NewAuthService(storage repositories.IAuthenticationRepository) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (a *AuthService) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	return a.storage.Register(req)
}

func (a *AuthService) LogIn(req models.LogInRequest) (*models.LoginResponse, error) {
	return a.storage.LogIn(req)
}

/*
Register(models.RegisterRequest) (*models.RegisterResponse, error)
LogIn(models.LogInRequest) (*models.LoginResponse, error)
*/
