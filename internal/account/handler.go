package account

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	List(ctx context.Context) ([]Account, error)
	GetByID(ctx context.Context, id int64) (Account, error)
	Create(ctx context.Context, account Account) (Account, error)
	Update(ctx context.Context, id int64, account Account) (Account, error)
}

type AccountHandler struct {
	svc Service
}

func NewAccountHandler(svc Service) *AccountHandler {
	return &AccountHandler{
		svc: svc,
	}
}

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	accounts, err := h.svc.List(c.Request.Context())
	if err != nil {
		slog.Error("failed to list accounts", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, toAccountResponses(accounts))
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, ok := parseAccountID(c)
	if !ok {
		return
	}

	account, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get account", "error", err, "id", id)
		switch {
		case errors.Is(err, ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, toAccountResponse(account))
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind create account request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	account, err := h.svc.Create(c.Request.Context(), toAccountFromRequest(req))
	if err != nil {
		slog.Error("failed to create account", "error", err)
		switch {
		case errors.Is(err, ErrInvalidAccount):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, toAccountResponse(account))
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, ok := parseAccountID(c)
	if !ok {
		return
	}

	var req AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind update account request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	account, err := h.svc.Update(c.Request.Context(), id, toAccountFromRequest(req))
	if err != nil {
		slog.Error("failed to update account", "error", err, "id", id)
		switch {
		case errors.Is(err, ErrInvalidAccount):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		case errors.Is(err, ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, toAccountResponse(account))
}

func parseAccountID(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("failed to parse account id", "error", err, "raw", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return 0, false
	}
	if id <= 0 {
		slog.Warn("non-positive account id", "id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return 0, false
	}

	return id, true
}
