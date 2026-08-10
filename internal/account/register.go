package account

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	repo := NewAccountRepository(db)
	svc := NewAccountService(repo)
	h := NewAccountHandler(svc)

	accounts := router.Group("/accounts")
	accounts.GET("/", h.ListAccounts)
	accounts.GET("/:id", h.GetAccount)
	accounts.POST("/", h.CreateAccount)
	accounts.PUT("/:id", h.UpdateAccount)
}
