package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/store"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	store *store.SQLiteStore
}

type CreatePostRequest struct {
	Title   string `json:"title" example:"First Post"`
	Content string `json:"content" example:"Hello from Echo and SQLite"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"post not found"`
}

func NewHandler(postStore *store.SQLiteStore) *Handler {
	return &Handler{store: postStore}
}

// ListPosts godoc
// @Summary List posts
// @Tags posts
// @Produce json
// @Success 200 {array} store.Post
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func (h *Handler) ListPosts(c *echo.Context) error {
	posts, err := h.store.ListPosts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to list posts"})
	}

	return c.JSON(http.StatusOK, posts)
}

// GetPost godoc
// @Summary Get post by ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} store.Post
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [get]
func (h *Handler) GetPost(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid post id"})
	}

	post, err := h.store.GetPost(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: "post not found"})
		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to fetch post"})
	}

	return c.JSON(http.StatusOK, post)
}

// CreatePost godoc
// @Summary Create a post
// @Tags posts
// @Accept json
// @Produce json
// @Param payload body CreatePostRequest true "Post payload"
// @Success 201 {object} store.Post
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts [post]
func (h *Handler) CreatePost(c *echo.Context) error {
	var request CreatePostRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
	}

	request.Title = strings.TrimSpace(request.Title)
	request.Content = strings.TrimSpace(request.Content)
	if request.Title == "" || request.Content == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "title and content are required"})
	}

	post, err := h.store.CreatePost(c.Request().Context(), store.CreatePostInput{
		Title:   request.Title,
		Content: request.Content,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to create post"})
	}

	return c.JSON(http.StatusCreated, post)
}
