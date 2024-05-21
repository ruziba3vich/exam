package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/exam/api/middleware"
	"github.com/ruziba3vich/exam/internal/models"
	"github.com/ruziba3vich/exam/internal/services"
)

type authHandler struct {
	service services.AuthService
	auth    middleware.Authentication
}

type AuthHandler struct {
	Service services.AuthService
	Auth    middleware.Authentication
}

func (a *AuthHandler) NewAuthHandler() *authHandler {
	return &authHandler{
		service: a.Service,
		auth:    a.Auth,
	}
}

func (a *authHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.Logger)
		return
	}
	response, err := a.service.Register(req)
	if err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.Logger)
		return
	}
	token, err := a.auth.GenerateToken(
		models.GenerateTokenRequest{
			Id:   response.Id,
			Name: response.Name,
		})
	if err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.Logger)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"token":    token,
		"response": response,
	})
}

func (a *authHandler) LogIn(c *gin.Context) {
	var req models.LogInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.Logger)
		return
	}
	response, err := a.service.LogIn(req)
	if err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.Logger)
		return
	}
	token, err := a.auth.GenerateToken(
		models.GenerateTokenRequest{
			Id:   response.Id,
			Name: response.Name,
		})
	if err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.Logger)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"token":    token,
		"response": response,
	})
}
