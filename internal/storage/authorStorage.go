package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ruziba3vich/exam/internal/models"
)

type Storage struct {
	db  *sql.DB
	ctx *context.Context
}

func New(db *sql.DB, ctx *context.Context) *Storage {
	return &Storage{
		db:  db,
		ctx: ctx,
	}
}

// / method to get all books
func (s *Storage) GetAllBooks(pagination, limit int) (results []models.Response, e error) {
	query := `
	SELECT
		b.id,
		b.title,
		a.name AS author_name,
		b.publication_date,
		b.isbn,
		b.description,
		b.created_at,
		b.updated_at
	FROM Books b
	INNER JOIN Authors a
	ON a.id = b.author_id
	WHERE b.is_deleted = false
	LIMIT $1 OFFSET $2;
	`
	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Microsecond*time.Duration(timeout))
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, limit, (pagination-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var response models.Response
		if err := rows.Scan(&response.Id,
			&response.Title,
			&response.AuthorName,
			&response.PublicationDate,
			&response.Isbn,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, response)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// / method to create new book
func (s *Storage) CreateBook(req models.CreateBookRequest) (*models.CreatedBookResponse, error) {
	query := `
		WITH inserted_book AS (
			INSERT INTO Books (
				title,
				description,
				author_id,
				isbn,
				created_at,
				updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, author_id
		)
		SELECT
			ib.id,
			ib.title,
			ib.description,
			a.name AS author_name,
			ib.isbn,
			ib.created_at,
			ib.updated_at
		FROM inserted_book ib
		JOIN Authors a ON ib.author_id = a.id;
	`

	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Microsecond*time.Duration(timeout))
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query,
		req.Title,
		req.Description,
		req.AuthorId,
		isbnGenerator(req.AuthorId),
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response models.CreatedBookResponse
	if err := rows.Scan(&response.Id,
		&response.Title,
		&response.Description,
		&response.AuthorName,
		&response.Isbn,
		&response.CreatedAt,
		&response.UpdatedAt); err != nil {
		return nil, err
	}
	return &response, nil
}

// / method to get book by id
func (s *Storage) GetBookById(id int) (*models.Response, error) {
	query := `
		SELECT
			b.id,
			b.title,
			a.name AS author_name,
			b.publication_date,
			b.isbn,
			b.description,
			b.created_at,
			b.updated_at
		FROM Books b
		INNER JOIN Authors a
		ON a.id = b.author_id
		WHERE b.is_deleted = false AND b.id = $1;
	`
	t, _ := strconv.Atoi((os.Getenv("TIMEOUT")))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(t))
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response models.Response
	if err := rows.Scan(&response.Id,
		&response.Title,
		&response.AuthorName,
		&response.PublicationDate,
		&response.Isbn,
		&response.Description,
		&response.CreatedAt,
		&response.UpdatedAt); err != nil {
		return nil, err
	}
	return &response, nil
}

// / method to update a book
func (s *Storage) UpdateBookById(req models.UpdateBookRequest) (*models.UpdatedBookResponse, error) {
	var query string
	t, _ := strconv.Atoi((os.Getenv("TIMEOUT")))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(t))
	defer cancel()
	if len(req.Title) > 0 && len(req.Description) > 0 && req.PublicationDate != nil {
		query = `
		UPDATE Books b
			INNER JOIN Authors a ON b.author_id = a.id
		SET b.title = $1,
			b.description = $2,
			b.publication_date = $3,
			b.updated_at = $4
		WHERE b.id = $5
		RETURNING
			b.id,
			b.title,
			b.description,
			a.name AS author_name,
			b.publication_date,
			b.created_at,
			b.updated_at;
		`

		rows, err := s.db.QueryContext(ctx,
			query,
			req.Title,
			req.Description,
			req.PublicationDate,
			time.Now(),
			req.GetBookId())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var response models.UpdatedBookResponse

		if err := rows.Scan(
			&response.Id,
			&response.Title,
			&response.Description,
			&response.AuthorName,
			&response.PublicationDate,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}

		return &response, nil
	} else if len(req.Title) > 0 && (len(req.Description) == 0 && req.PublicationDate == nil) {
		query = `
			UPDATE Books b
				INNER JOIN Authors a ON b.author_id = a.id
			SET b.title = $1
			SET b.updated_at = $2
			WHERE b.id = $3
			RETURNING
				b.id,
				b.title,
				b.description,
				a.name AS author_name,
				b.publication_date,
				b.created_at,
				b.updated_at;
			`

		rows, err := s.db.QueryContext(ctx,
			query,
			req.Title,
			time.Now(),
			req.GetBookId())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var response models.UpdatedBookResponse

		if err := rows.Scan(
			&response.Id,
			&response.Title,
			&response.Description,
			&response.AuthorName,
			&response.PublicationDate,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}

		return &response, nil
	} else if len(req.Description) > 0 && (len(req.Title) == 0 && req.PublicationDate == nil) {
		query = `
			UPDATE Books b
				SET b.description = $1
				SET b.updated_at = $2
			WHERE b.id = $3;
		`
		rows, err := s.db.QueryContext(ctx,
			query,
			req.Description,
			time.Now(),
			req.GetBookId())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var response models.UpdatedBookResponse

		if err := rows.Scan(
			&response.Id,
			&response.Title,
			&response.Description,
			&response.AuthorName,
			&response.PublicationDate,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}

		return &response, nil
	} else if req.PublicationDate == nil && (len(req.Title) == 0 && len(req.Description) == 0) {
		query = `
			UPDATE Books b
				SET b.publication_date = $1
				SET b.updated_at = $2
			WHERE b.id = $3;
		`
		rows, err := s.db.QueryContext(ctx,
			query,
			req.PublicationDate,
			time.Now(),
			req.GetBookId())
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var response models.UpdatedBookResponse

		if err := rows.Scan(
			&response.Id,
			&response.Title,
			&response.Description,
			&response.AuthorName,
			&response.PublicationDate,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}

		return &response, nil
	} else {
		return nil, errors.New("invalid request")
	}
}

// / method to delete a book
func (s *Storage) DeleteBookById(bookId int) (*models.DeletedBookResponse, error) {
	t, _ := strconv.Atoi((os.Getenv("TIMEOUT")))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(t))
	defer cancel()
	query := `
		UPDATE Books
		SET is_deleted = true
		WHERE id = $1
		RETURNING
			b.id,
			b.title,
			b.description,
			a.name AS author_name,
			b.publication_date,
			b.created_at,
			b.updated_at;
	`

	rows, err := s.db.QueryContext(ctx, query, bookId)
	if err != nil {
		return nil, err
	}
	var response models.DeletedBookResponse

	if err := rows.Scan(
		&response.Id,
		&response.Title,
		&response.Description,
		&response.AuthorName,
		&response.PublicationDate,
		&response.CreatedAt,
		&response.UpdatedAt); err != nil {
		return nil, err
	}

	return &response, nil
}

func isbnGenerator(creatorID int) string {
	prefix := "978"

	registrationGroup := "0"

	registrantElement := fmt.Sprintf("%07d", creatorID%10000000)

	publicationElement := fmt.Sprintf("%06d", creatorID%1000000)

	partialISBN := strings.Join([]string{prefix, registrationGroup, registrantElement, publicationElement}, "")

	checkDigit := calculateCheckDigit(partialISBN)

	isbn := partialISBN + strconv.Itoa(checkDigit)
	return isbn
}

func calculateCheckDigit(isbn string) int {
	sum := 0
	for i, digit := range isbn {
		num := int(digit - '0')
		if i%2 == 0 {
			sum += num
		} else {
			sum += num * 3
		}
	}
	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

/*
type Book struct {
	id              int       //
	Title           string    `json:"title"`
	authorId        int       //
	PublicationDate time.Time `json:"publication_date"`
	isbn            string    //
	Description     string    `json:"description"`
	createdAt       time.Time //
	updatedAt       time.Time //
	isDeleted       bool
}
*/

/*
	/// CreateBook(models.CreateBookRequest) (models.CreatedBookResponse, error)
	/// GetAllBooks() ([]models.Response, error)
	/// GetBookById(int) (models.Response, error)
	/// UpdateBookById(models.UpdateBookRequest) (models.UpdatedBookResponse, error)
	DeleteBookById(int) (models.DeletedBookResponse, error)
	GetBooksByAuthor(int) ()
*/
