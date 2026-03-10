package main

import (
	"log"
	"os"

	_ "backend/docs"
	"backend/internal/api"
	"backend/internal/store"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// @title Posts API
// @version 1.0
// @description Simple Echo API with a local Turso database file.
// @BasePath /
// @schemes http
func main() {
	dbPath := envOrDefault("DB_PATH", "app.db")
	port := envOrDefault("PORT", "8080")

	postStore, err := store.NewStore(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize Turso store: %v", err)
	}
	defer postStore.Close()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	handler := api.NewHandler(postStore)

	e.GET("/posts", handler.ListPosts)
	e.GET("/posts/:id", handler.GetPost)
	e.POST("/posts", handler.CreatePost)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Printf("server listening on :%s using Turso database file %q", port, dbPath)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
