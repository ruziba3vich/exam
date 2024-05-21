package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/exam/api/middleware"
	"github.com/ruziba3vich/exam/internal/models"
	"github.com/ruziba3vich/exam/internal/repositories"
)

type handler struct {
	storage repositories.IAuthorRepository
	auth    middleware.Authentication
}

type HandlerConfig struct {
	Storage repositories.IAuthorRepository
	Auth    middleware.Authentication
}

func NewHandler(h *HandlerConfig) *handler {
	return &handler{
		storage: h.Storage,
		auth:    h.Auth,
	}
}

// / handler to create a new book
func (h *handler) CreateBookHandler(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.Logger)
		return
	}
	var request models.CreateBookRequest
	id, err := h.auth.ExtractIdFromToken(c)
	if err != nil {
		printError(http.StatusUnauthorized, err, c, h.auth.Logger)
		return
	}
	request.SetId(id)
	name, err := h.auth.ExtractAuthorNameFromToken(c)
	if err != nil {
		printError(http.StatusUnauthorized, err, c, h.auth.Logger)
		return
	}
	request.SetName(name)
	if err := c.ShouldBindJSON(&request); err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	result, err := h.storage.CreateBook(request)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

var defaultPage int
var defaultLimit int = 10

// / handler to get all books which belong to author
func (h *handler) GetAllBooksHandler(c *gin.Context) {
	pagination := c.Query("page")
	if pagination == "" {
		pagination = fmt.Sprintf("%d", defaultPage)
		defaultPage++
	}
	limit := c.Query("limit")
	if limit == "" {
		limit = fmt.Sprintf("%d", defaultLimit)
		defaultLimit += 10
	}

	pg, _ := strconv.Atoi(pagination)
	lm, _ := strconv.Atoi(limit)
	books, err := h.storage.GetAllBooks(pg, lm)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

// / handler to get a book by id
func (h *handler) GetBookById(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.Logger)
		return
	}
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}

	book, err := h.storage.GetBookById(id)
	if err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// / handler to update a book by id
func (h *handler) UpdateBookById(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.Logger)
		return
	}
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	var request models.UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	request.SetBookId(id)
	newBook, err := h.storage.UpdateBookById(request)
	if err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	c.IndentedJSON(http.StatusOK, newBook)
}

// / handler to delete a book and return it
func (h *handler) DeleteBookById(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.Logger)
		return
	}
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	deletedBook, err := h.storage.DeleteBookById(id)
	if err != nil {
		h.auth.Logger.Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.Logger)
		return
	}
	c.IndentedJSON(http.StatusOK, deletedBook)
}

// / error printer function
func printError(status int, err error, c *gin.Context, lg *log.Logger) {
	c.IndentedJSON(status, gin.H{"error": err})
	lg.Println(err)
}

/*
	CreateBook(models.CreateBookRequest) (models.CreatedBookResponse, error)
    GetAllBooks(int) ([]models.Response, error)
    GetBookById(int) (models.Response, error)
    UpdateBookById(models.UpdateBookRequest) (models.UpdatedBookResponse, error)
    DeleteBookById(int) (models.DeletedBookResponse, error)
*/
