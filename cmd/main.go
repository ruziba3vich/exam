package main

import (
	"context"
	"log"
	"os"

	"github.com/ruziba3vich/exam/api"
	"github.com/ruziba3vich/exam/api/middleware"
	"github.com/ruziba3vich/exam/internal/storage"
)

func main() {
	db := storage.DB()
	ctx := context.Background()
	logger := log.New(os.Stdout, "app : ", log.Flags())

	auth := storage.NewAuthService(db, &ctx)
	storage := storage.New(db, &ctx)
	middleware := middleware.New(logger)

	server := api.New(
		api.API{
			WorkerService: storage,
			AuthService:   auth,
			MiddleWare:    middleware,
		},
	)

	if err := server.Run(":7777"); err != nil {
		log.Fatal("Failed to run HTTP server:  ", err)
		panic(err)
	}

	log.Print("Server stopped")
}
