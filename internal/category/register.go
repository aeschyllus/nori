package category

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	repo := NewCategoryRepository(db)
	svc := NewCategoryService(repo)
	h := NewCategoryHandler(svc)

	categories := router.Group("/categories")
	categories.GET("/", h.ListCategories)
	categories.GET("/:id", h.GetCategory)
	categories.POST("/", h.CreateCategory)
	categories.PUT("/:id", h.UpdateCategory)
}
