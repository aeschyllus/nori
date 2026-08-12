package category

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	List(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id int64) (Category, error)
	Create(ctx context.Context, category Category) (Category, error)
	Update(ctx context.Context, id int64, category Category) (Category, error)
}

type CategoryHandler struct {
	svc Service
}

func NewCategoryHandler(svc Service) *CategoryHandler {
	return &CategoryHandler{
		svc: svc,
	}
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.svc.List(c.Request.Context())
	if err != nil {
		slog.Error("failed to list categories", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toCategoryResponses(categories))
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, ok := parseCategoryID(c)
	if !ok {
		return
	}

	category, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get category", "error", err, "id", id)
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, toCategoryResponse(category))
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind create category request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	category, err := h.svc.Create(c.Request.Context(), toCategoryFromRequest(req))
	if err != nil {
		slog.Error("failed to create category", "error", err)
		switch {
		case errors.Is(err, ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, toCategoryResponse(category))
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, ok := parseCategoryID(c)
	if !ok {
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind update category request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	category, err := h.svc.Update(c.Request.Context(), id, toCategoryFromRequest(req))
	if err != nil {
		slog.Error("failed to update category", "error", err, "id", id)
		switch {
		case errors.Is(err, ErrInvalidCategory):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		case errors.Is(err, ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, toCategoryResponse(category))
}

func parseCategoryID(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("failed to parse category id", "error", err, "raw", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return 0, false
	}
	if id <= 0 {
		slog.Warn("non-positive category id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return 0, false
	}

	return id, true
}
