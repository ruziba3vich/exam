package repositories

import "github.com/ruziba3vich/exam/internal/models"

type IAuthorRepository interface {
	CreateBook(models.CreateBookRequest) (*models.CreatedBookResponse, error)
	GetAllBooks(int, int) ([]models.Response, error)
	GetBookById(int) (*models.GetBookByIdResponse, error)
	UpdateBookById(models.UpdateBookRequest) (*models.UpdatedBookResponse, error)
	DeleteBookById(int) (*models.DeletedBookResponse, error)
	// GetBooksByAuthor(int) ([]models.Request, error)
	UpdateBiography(models.UpdateBiographyRequest) (*models.UpdateBiographyResponse, error)
	UpdateBirthdate(models.UpdateBirthdateRequest) (*models.UpdateBirthdateResponse, error)
	GetProfile(models.GetProfileRequest) (*models.GetProfileResponse, error)
	GetAllAuthors() ([]models.GetProfileResponse, error)
	GetAuthorById(id uint) (*models.GetProfileResponse, error)
}
