package category

type CategoryResponse struct {
	ID   int64  `json:"categoryId"`
	Name string `json:"name"`
}

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
