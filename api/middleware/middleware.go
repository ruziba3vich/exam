package middleware

import (
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	jwttokens "github.com/ruziba3vich/authentication_tokens"
	"github.com/ruziba3vich/exam/internal/models"
)

type Authentication struct {
	Logger *log.Logger
}

func New(logger *log.Logger) *Authentication {
	return &Authentication{
		Logger: logger,
	}
}

func (a *Authentication) GetLogger() *log.Logger {
	return a.Logger
}

func (a *Authentication) ValidateToken(c *gin.Context) (bool, error) {
	claims := &jwttokens.Claims{}

	tokenString, err := extractToken(c)
	if err != nil {
		return false, err
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwttokens.GetJwtToken(), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, fmt.Errorf("invalid token signature")
		}
		return false, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func (a *Authentication) GenerateToken(req models.GenerateTokenRequest) (string, error) {
	return jwttokens.GenerateToken(req.Id, req.Name)
}

func (a *Authentication) ExtractAuthorIdFromToken(c *gin.Context) (int, error) {
	token, err := extractToken(c)
	if err != nil {
		return 0, err
	}
	return jwttokens.ExtractUserIDFromToken(token)
}

func (a *Authentication) ExtractAuthorNameFromToken(c *gin.Context) (string, error) {
	token, err := extractToken(c)
	if err != nil {
		return token, err
	}
	return jwttokens.ExtractUsernameFromToken(token)
}

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}
	return "", errors.New("invalid token")
}
