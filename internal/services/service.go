package services

import (
	"github.com/ruziba3vich/exam/internal/models"
	"github.com/ruziba3vich/exam/internal/repositories"
)

type Service struct {
	storage repositories.IAuthorRepository
}

func New(storage repositories.IAuthorRepository) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateBook(req models.CreateBookRequest) (*models.CreatedBookResponse, error) {
	return s.storage.CreateBook(req)
}

func (s *Service) GetAllBooks(pagination, limit int) ([]models.Response, error) {
	return s.storage.GetAllBooks(pagination, limit)
}

func (s *Service) GetBookById(bookId int) (*models.GetBookByIdResponse, error) {
	return s.storage.GetBookById(bookId)
}

func (s *Service) UpdateBookById(req models.UpdateBookRequest) (*models.UpdatedBookResponse, error) {
	return s.storage.UpdateBookById(req)
}

func (s *Service) DeleteBookById(bookId int) (*models.DeletedBookResponse, error) {
	return s.storage.DeleteBookById(bookId)
}

func (s *Service) UpdateBiography(req models.UpdateBiographyRequest) (*models.UpdateBiographyResponse, error) {
	return s.storage.UpdateBiography(req)
}

func (s *Service) UpdateBirthdate(req models.UpdateBirthdateRequest) (*models.UpdateBirthdateResponse, error) {
	return s.storage.UpdateBirthdate(req)
}

/*
	CreateBook(models.CreateBookRequest) (models.CreatedBookResponse, error)
	GetAllBooks() ([]models.Response, error)
	GetBookById() (models.Response, error)
	UpdateBookById(models.UpdateBookRequest) (models.UpdatedBookResponse, error)
	DeleteBookById() (models.DeletedBookResponse, error)
*/
