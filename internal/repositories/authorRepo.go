package repositories

import "github.com/ruziba3vich/exam/internal/models"

type IAuthorRepository interface {
	CreateBook(models.CreateBookRequest) (*models.CreatedBookResponse, error)
	GetAllBooks(int, int) ([]models.Response, error)
	GetBookById(int) (*models.Response, error)
	UpdateBookById(models.UpdateBookRequest) (*models.UpdatedBookResponse, error)
	DeleteBookById(int) (*models.DeletedBookResponse, error)
	// GetBooksByAuthor(int) ([]models.Request, error)
}
