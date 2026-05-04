package main

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/topi0247/hitori/handler"
	"github.com/topi0247/hitori/infra/store"
	"github.com/topi0247/hitori/usecase"
)

func main() {
	ctx := context.Background()

	pool, err := store.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	repos := store.NewRepositories(pool)
	usecases := usecase.NewUsecases(repos)
	handlers := handler.NewHandlers(usecases)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	handlers.SetRoutes(e, os.Getenv("JWT_SECRET"))

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
