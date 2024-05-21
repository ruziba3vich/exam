package repositories

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/exam/internal/models"
)

type IMiddleWare interface {
	ValidateToken(*gin.Context) (bool, error)
	GenerateToken(models.GenerateTokenRequest) (string, error)
	ExtractAuthorIdFromToken(c *gin.Context) (int, error)
	ExtractAuthorNameFromToken(c *gin.Context) (string, error)
	GetLogger() *log.Logger
}
