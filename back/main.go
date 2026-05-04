package main

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/topi0247/hitori/handler"
	"github.com/topi0247/hitori/infra/postgres"
	"github.com/topi0247/hitori/usecase"
)

func main() {
	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	repos := postgres.NewRepositories(pool)
	usecases := usecase.NewUsecases(repos)
	handlers := handler.NewHandlers(usecases)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	handlers.SetRoutes(e)

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
