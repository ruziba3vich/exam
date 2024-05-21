package storage

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/ruziba3vich/exam/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImple struct {
	db  *sql.DB
	ctx *context.Context
}

func NewAuthService(db *sql.DB, ctx *context.Context) *AuthServiceImple {
	return &AuthServiceImple{
		db:  db,
		ctx: ctx,
	}
}

func (a *AuthServiceImple) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	ctx, cancel := getContext(a.ctx)
	defer cancel()
	hashshedPwd, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO Authors (
			name,
			password,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4
		)
		RETURNING id, name, password, created_at, updated_at;
	`
	rows, err := a.db.QueryContext(ctx, query, req.Name, hashshedPwd, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var response models.RegisterResponse
	if err := rows.Scan(
		&response.Id,
		&response.Name,
		&response.Password,
		&response.CreatedAt,
		&response.UpdatedAt); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &response, nil
}

func (a *AuthServiceImple) LogIn(req models.LogInRequest) (*models.LoginResponse, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getContext(a.ctx)
	defer cancel()
	query := `
		SELECT
			a.id,
			a.name,
			a.password,
			a.biography,
			a.birth_date,
			a.created_at,
			a.updated_at
		FROM Authors a
		WHERE a.name = $1 AND a.password = $2;
	`

	rows, err := a.db.QueryContext(ctx, query, req.Name, hashedPassword)
	if err != nil {
		return nil, err
	}
	var response models.LoginResponse
	if err := rows.Scan(
		&response.Id,
		&response.Name,
		&response.Password,
		&response.Biography,
		&response.CreatedAt,
		&response.UpdatedAt); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &response, nil
}

func hashPassword(password string) (string, error) {
	cost := 111

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func getContext(c *context.Context) (context.Context, context.CancelFunc) {
	t, _ := strconv.Atoi((os.Getenv("TIMEOUT")))
	return context.WithTimeout(*c, time.Millisecond*time.Duration(t))
}

/*
type Author struct {
	id        int
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	BirthDate *time.Time `json:"birth_date"`
	Biography string     `json:"biography"`
	createdAt time.Time
	updatedAt time.Time
}
*/
