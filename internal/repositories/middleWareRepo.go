package repositories

import "github.com/ruziba3vich/exam/internal/models"

type IMiddleWare interface {
	GenerateToken(models.GenerateTokenRequest) models.GenerateTokenResponse
	ExtractAuthorIdFromToken(models.ExtractAuthorIdFromTokenRequest) models.ExtractIdFromTokenResponse
	ExtractAuthorNameFromToken(models.ExtractAuthorNameFromTokenRequest) models.ExtractAuthorNameFromTokenResponse
}
