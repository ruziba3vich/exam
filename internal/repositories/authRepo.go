package repositories

import "github.com/ruziba3vich/exam/internal/models"

type IAuthenticationRepository interface {
	Register(models.RegisterRequest) (*models.RegisterResponse, error)
	LogIn(models.LogInRequest) (*models.LoginResponse, error)
}
