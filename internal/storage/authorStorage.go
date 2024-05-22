package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/pp"
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
func (s *Storage) GetAllBooks(pagination, limit int) (results []models.Response, err error) {
	query := `
	SELECT
		b.id,
		b.title,
		a.id,
		b.publication_date,
		b.isbn,
		b.description,
		b.created_at,
		b.updated_at
	FROM Books b
	INNER JOIN Authors a
	ON a.id = b.author_id
	WHERE b.is_deleted = false
	`

	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))

	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(timeout))
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	pp.Println("selected all")
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var response models.Response
		err := rows.Scan(&response.Id,
			&response.Title,
			&response.AuthorId,
			&response.PublicationDate,
			&response.Isbn,
			&response.Description,
			&response.CreatedAt,
			&response.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		results = append(results, response)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}

// / method to create new book
func (s *Storage) CreateBook(req models.CreateBookRequest) (*models.CreatedBookResponse, error) {
	pp.Println(req, "\n---------------came-----------------------------")

	query := `
	INSERT INTO Books (
		title,
		description,
		author_id,
		isbn,
		publication_date,
		created_at,
		updated_at,
		is_deleted
	)
	VALUES ($1, $2, $3, $4, COALESCE($5, NOW()), $6, $7, $8)
	RETURNING id, title, description, author_id, isbn, publication_date, created_at, updated_at;
	`

	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Microsecond*time.Duration(timeout))
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query,
		req.Title,
		req.Description,
		req.AuthorId,
		isbnGenerator(req.AuthorId),
		req.PublicationDate,
		time.Now(),
		time.Now(),
		false)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var response models.CreatedBookResponse
	for rows.Next() {
		if err := rows.Scan(
			&response.Id,
			&response.Title,
			&response.Description,
			&response.AuthorId,
			&response.Isbn,
			&response.PublicationDate,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}
	}
	return &response, nil
}

// / method to get book by id
func (s *Storage) GetBookById(id int) (*models.GetBookByIdResponse, error) {
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
	t, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(t))
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response models.GetBookByIdResponse
	for rows.Next() {
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
	}
	return &response, nil
}

// / method to update a book
func (s *Storage) UpdateBookById(req models.UpdateBookRequest) (*models.UpdatedBookResponse, error) {
	var query string
	t, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(t))
	defer cancel()
	/// update whole book
	if len(req.Title) > 0 && len(req.Description) > 0 && req.PublicationDate != nil {
		pp.Println("1111111111111111111111---111111")
		query = `
		UPDATE Books
		SET 
			title = $1,
			description = $2,
			publication_date = $3,
			updated_at = $4
		WHERE 
			id = $5
		RETURNING 
			id,
			title,
			description,
			author_id,
			publication_date,
			created_at,
			updated_at;
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

		for rows.Next() {
			if err := rows.Scan(
				&response.Id,
				&response.Title,
				&response.Description,
				&response.AuthorId,
				&response.PublicationDate,
				&response.CreatedAt,
				&response.UpdatedAt); err != nil {
				return nil, err
			}
		}

		return &response, nil

		/// update title only | done
	} else if len(req.Title) > 0 && (len(req.Description) == 0 && req.PublicationDate == nil) {
		pp.Println("222222222222222222222222222---2222222222")
		query = `
			UPDATE Books
			SET title = $1,
				updated_at = $2
			WHERE id = $3
			RETURNING
				id,
				title,
				description,
				author_id,
				publication_date,
				created_at,
				updated_at;
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

		for rows.Next() {
			if err := rows.Scan(
				&response.Id,
				&response.Title,
				&response.Description,
				&response.AuthorId,
				&response.PublicationDate,
				&response.CreatedAt,
				&response.UpdatedAt); err != nil {
				return nil, err
			}
		}

		return &response, nil
		/// update description | done
	} else if len(req.Description) > 0 && (len(req.Title) == 0 && req.PublicationDate == nil) {
		pp.Println("3333333333333333333333333---333333333333333")
		query = `
			UPDATE Books
				SET description = $1,
					updated_at = $2
			WHERE id = $3
			RETURNING
				id,
				title,
				description,
				author_id,
				publication_date,
				created_at,
				updated_at;;
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

		for rows.Next() {
			if err := rows.Scan(
				&response.Id,
				&response.Title,
				&response.Description,
				&response.AuthorId,
				&response.PublicationDate,
				&response.CreatedAt,
				&response.UpdatedAt); err != nil {
				return nil, err
			}
		}

		return &response, nil

		/// udate publication date only
	} else if req.PublicationDate != nil && (len(req.Title) == 0 && len(req.Description) == 0) {
		pp.Println("444444444444444444444444444444---444444444")
		query = `
			UPDATE Books
				SET publication_date = $1,
					updated_at = $2
			WHERE id = $3
			RETURNING
				id,
				title,
				description,
				author_id,
				publication_date,
				created_at,
				updated_at;;
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

		for rows.Next() {
			if err := rows.Scan(
				&response.Id,
				&response.Title,
				&response.Description,
				&response.AuthorId,
				&response.PublicationDate,
				&response.CreatedAt,
				&response.UpdatedAt); err != nil {
				return nil, err
			}
		}

		return &response, nil
	} else {
		pp.Println("555555555555555555555---5555555555555555555555")
		return nil, errors.New("invalid request")
	}
}

// / method to delete a book
func (s *Storage) DeleteBookById(bookId int) (*models.DeletedBookResponse, error) {
	t, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	ctx, cancel := context.WithTimeout(*s.ctx, time.Millisecond*time.Duration(t))
	defer cancel()
	query := `
		UPDATE Books
		SET is_deleted = true
		WHERE id = $1
		RETURNING
			id,
			title,
			description,
			author_id,
			publication_date,
			created_at,
			updated_at;
	`

	rows, err := s.db.QueryContext(ctx, query, bookId)
	if err != nil {
		return nil, err
	}
	var response models.DeletedBookResponse

	for rows.Next() {
		if err := rows.Scan(
			&response.Id,
			&response.Title,
			&response.Description,
			&response.AuthorId,
			&response.PublicationDate,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			return nil, err
		}
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
	log.Println(isbn, "-------------generated-------------------")
	return isbn
}

func calculateCheckDigit(isbn string) int {
	log.Println(isbn, "--------------------------------")
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
	log.Println(checkDigit, "------------check digit--------------------")
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
