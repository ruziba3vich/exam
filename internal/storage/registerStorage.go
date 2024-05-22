package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ruziba3vich/exam/internal/models"
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
	hashshedPwd := hashPassword(req.Password)
	// if err != nil {
	// 	return nil, err
	// }
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
		log.Println(err, "eeeeeeeeeeeeeeeeeeerrrrrrrrrrroooooooooor")
		return nil, err
	}
	defer rows.Close()
	var response models.RegisterResponse
	for rows.Next() {
		if err := rows.Scan(
			&response.Id,
			&response.Name,
			&response.Password,
			&response.CreatedAt,
			&response.UpdatedAt); err != nil {
			log.Println("5555555555555555555", err)
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("6666666666666666666")
		return nil, err
	}
	return &response, nil
}

func (a *AuthServiceImple) LogIn(req models.LogInRequest) (*models.LoginResponse, error) {
	hashedPassword := hashPassword(req.Password)
	// if err != nil {
	// 	return nil, err
	// }
	ctx, cancel := getContext(a.ctx)
	defer cancel()
	query := `
		SELECT
			a.id,
			a.name,
			a.password
		FROM Authors a
		WHERE a.name = $1 AND a.password = $2;
	`

	rows, err := a.db.QueryContext(ctx, query, req.Name, hashedPassword)
	if err != nil {
		log.Println("7777777777777777777777777", err)
		return nil, err
	}
	var response models.LoginResponse
	for rows.Next() {
		if err := rows.Scan(
			&response.Id,
			&response.Name,
			&response.Password); err != nil {
			log.Println("8888888888888888888888888888888", err)
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("99999999999999999999999999999999", err)
		return nil, err
	}
	return &response, nil
}

func hashPassword(password string) string {
	hasher := sha256.New()

	hasher.Write([]byte(password))

	hashedPasswordBytes := hasher.Sum(nil)

	hashedPassword := hex.EncodeToString(hashedPasswordBytes)

	return hashedPassword
}

func getContext(c *context.Context) (context.Context, context.CancelFunc) {
	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		log.Println(timeoutStr, "4555-----------------empty------------------------------")
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println(timeoutStr, "4333-----------------invalid value------------------------------")
	}
	return context.WithTimeout(*c, time.Millisecond*time.Duration(timeout))
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
