package models

import (
	"context"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LogInRequest struct {
	RegisterRequest
}

//////////////

type Request struct {
	id              int
	bookId          int
	name            string
	Title           string     `json:"title"`
	AuthorId        int        `json:"author_id"`
	Description     string     `json:"description"`
	PublicationDate *time.Time `json:"publication_date"`
	ctx             *context.Context
}

func (r *Request) SetContext(context *context.Context) {
	r.ctx = context
}

func (r *Request) GetContext() *context.Context {
	return r.ctx
}

type CreateBookRequest struct {
	Request
}

type UpdateBookRequest struct {
	Request
}

type GenerateTokenRequest struct {
	Id   int
	Name string
}

type Token struct {
	Token string
}

type ExtractAuthorNameFromTokenRequest struct {
	Token
}

type ExtractAuthorIdFromTokenRequest struct {
	Token
}

func (r *Request) SetId(id int) {
	r.id = id
}

func (r *Request) SetName(name string) {
	r.name = name
}

func (r *Request) GetId() int {
	return r.id
}

func (r *Request) GetName() string {
	return r.name
}

func (r *Request) SetBookId(id int) {
	r.bookId = id
}

func (r *Request) GetBookId() int {
	return r.bookId
}
