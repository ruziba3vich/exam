package models

import "time"

// ////////////////////////////////////////////////////////////////////////////////////////////////
type RegisterResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Biography string    `json:"biography"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	RegisterResponse
}

// ////////////////////////////////////////////////////////////////////////////////////////////////

type Response struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	AuthorName      string    `json:"author_name"`
	PublicationDate time.Time `json:"publication_date"`
	Isbn            string    `json:"isbn"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreatedBookResponse struct {
	Response
}

type UpdatedBookResponse struct {
	Response
}

type DeletedBookResponse struct {
	Response
}

type GenerateTokenResponse struct {
	Token string
	Error error
}

type ExtractIdFromTokenResponse struct {
	Id    int
	Error error
}

type ExtractAuthorNameFromTokenResponse struct {
	Name  string
	Error error
}
