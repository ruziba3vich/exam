package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/exam/api/handlers"
	"github.com/ruziba3vich/exam/internal/repositories"
)

type API struct {
	WorkerService repositories.IAuthorRepository
	AuthService   repositories.IAuthenticationRepository
	MiddleWare    repositories.IMiddleWare
}

func New(api API) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handler := handlers.NewHandler(
		&handlers.HandlerConfig{
			Storage:     api.WorkerService,
			Auth:        api.MiddleWare,
			AuthService: api.AuthService,
		},
	)

	router.GET("/", handler.GetProfile)

	router.POST("/create/book", handler.CreateBookHandler)
	router.POST("/register", handler.Register)
	router.POST("/login", handler.LogIn)
	router.GET("/authors", handler.GetAllAuthors)
	router.GET("/authors/:id", handler.GetAuthorById)
	router.GET("/books/", handler.GetAllBooksHandler)
	router.GET("/books/:id", handler.GetBookById)
	router.PUT("/update/book/:id", handler.UpdateBookById)
	router.PUT("/update/birthdate/:id", handler.UpdateBirthdate)
	router.PUT("/update/biography/:id", handler.UpdateBiography)
	router.DELETE("/delete/book/:id", handler.DeleteBookById)

	return router
}
