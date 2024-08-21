package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/exam/api/handlers"
	"github.com/ruziba3vich/exam/api/metrics"
	"github.com/ruziba3vich/exam/internal/repositories"
)

type API struct {
	WorkerService repositories.IAuthorRepository
	AuthService   repositories.IAuthenticationRepository
	MiddleWare    repositories.IMiddleWare
}

func New(api API) *gin.Engine {
	metric := metrics.New()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		metric.Requests.Add(1)

		c.Next()

		status := c.Writer.Status()
		if status >= 200 && status < 400 {
			metric.Successes.Add(1)
		} else {
			metric.Errors.Add(1)
		}
	})

	handler := handlers.NewHandler(
		&handlers.HandlerConfig{
			Storage:     api.WorkerService,
			Auth:        api.MiddleWare,
			AuthService: api.AuthService,
		},
	)

	router.GET("/", handler.GetProfile)

	router.POST("/create/book", metric.IncreaseDbRequests, handler.CreateBookHandler)
	router.POST("/register", metric.IncreaseDbRequests, handler.Register)
	router.POST("/login", metric.IncreaseDbRequests, handler.LogIn)
	router.GET("/authors", handler.GetAllAuthors)
	router.GET("/authors/:id", handler.GetAuthorById)
	router.GET("/books/", handler.GetAllBooksHandler)
	router.GET("/books/:id", handler.GetBookById)
	router.PUT("/update/book/:id", metric.IncreaseDbRequests, handler.UpdateBookById)
	router.PUT("/update/birthdate/:id", metric.IncreaseDbRequests, handler.UpdateBirthdate)
	router.PUT("/update/biography/:id", metric.IncreaseDbRequests, handler.UpdateBiography)
	router.DELETE("/delete/book/:id", metric.IncreaseDbRequests, handler.DeleteBookById)

	router.GET("/debug/vars", gin.WrapH(http.DefaultServeMux))

	return router
}
