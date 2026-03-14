package main

import (
	"backend/internal/api"
	"backend/internal/store"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/syumai/workers"
)

// @title Posts API
// @version 1.0
// @description Simple Echo API with SQLite-compatible storage.
// @BasePath /
// @schemes http
func main() {
	postStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}
	defer postStore.Close()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins(),
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	handler := api.NewHandler(postStore)

	e.GET("/posts", handler.ListPosts)
	e.GET("/posts/:id", handler.GetPost)
	e.POST("/posts", handler.CreatePost)

	log.Printf("serving posts API using %s", store.SourceDescription())
	workers.Serve(e)
}

func allowedOrigins() []string {
	value := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
	if value == "" {
		return []string{"*"}
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return []string{"*"}
	}

	return origins
}
