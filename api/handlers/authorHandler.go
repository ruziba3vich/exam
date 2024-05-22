package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/exam/internal/models"
	"github.com/ruziba3vich/exam/internal/repositories"
)

type handler struct {
	storage     repositories.IAuthorRepository
	authService repositories.IAuthenticationRepository
	auth        repositories.IMiddleWare
}

type HandlerConfig struct {
	Storage     repositories.IAuthorRepository
	AuthService repositories.IAuthenticationRepository
	Auth        repositories.IMiddleWare
}

func NewHandler(h *HandlerConfig) *handler {
	return &handler{
		storage:     h.Storage,
		authService: h.AuthService,
		auth:        h.Auth,
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////
func (a *handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("11111111111111111111")
		printError(http.StatusBadRequest, err, c, a.auth.GetLogger())
		return
	}
	response, err := a.authService.Register(req)
	if err != nil {
		log.Println("222222222222222222222222")
		printError(http.StatusBadRequest, err, c, a.auth.GetLogger())
		return
	}
	token, err := a.auth.GenerateToken(
		models.GenerateTokenRequest{
			Id:   response.Id,
			Name: response.Name,
		})
	if err != nil {
		log.Println("333333333333333333333333")
		printError(http.StatusBadRequest, err, c, a.auth.GetLogger())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"token":    token,
		"response": response,
	})
}

func (a *handler) LogIn(c *gin.Context) {
	var req models.LogInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.GetLogger())
		return
	}
	response, err := a.authService.LogIn(req)
	if err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.GetLogger())
		return
	}
	token, err := a.auth.GenerateToken(
		models.GenerateTokenRequest{
			Id:   response.Id,
			Name: response.Name,
		})
	if err != nil {
		printError(http.StatusBadRequest, err, c, a.auth.GetLogger())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"token":    token,
		"response": response,
	})
}

////////////////////////////////////////////////////////////////////////////////////////////

var defaultPage int
var defaultLimit int = 10

// / handler to create a new book
func (h *handler) CreateBookHandler(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.GetLogger())
		return
	}
	var request models.CreateBookRequest
	id, err := h.auth.ExtractAuthorIdFromToken(c)
	if err != nil {
		printError(http.StatusUnauthorized, err, c, h.auth.GetLogger())
		return
	}
	request.AuthorId = id
	name, err := h.auth.ExtractAuthorNameFromToken(c)
	if err != nil {
		printError(http.StatusUnauthorized, err, c, h.auth.GetLogger())
		return
	}
	request.SetName(name)
	if err := c.ShouldBindJSON(&request); err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	result, err := h.storage.CreateBook(request)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

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
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

// / handler to get a book by id
func (h *handler) GetBookById(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.GetLogger())
		return
	}
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		h.auth.GetLogger().Println(err)
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}

	book, err := h.storage.GetBookById(id)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// / handler to update a book by id
func (h *handler) UpdateBookById(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.GetLogger())
		return
	}
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	var request models.UpdateBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	request.SetBookId(id)
	newBook, err := h.storage.UpdateBookById(request)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	c.IndentedJSON(http.StatusOK, newBook)
}

// / handler to delete a book and return it
func (h *handler) DeleteBookById(c *gin.Context) {
	ok, err := h.auth.ValidateToken(c)
	if !ok {
		printError(http.StatusForbidden, err, c, h.auth.GetLogger())
		return
	}
	bookId := c.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
		return
	}
	deletedBook, err := h.storage.DeleteBookById(id)
	if err != nil {
		printError(http.StatusBadRequest, err, c, h.auth.GetLogger())
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
