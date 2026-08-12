package category

func toCategoryResponse(c Category) CategoryResponse {
	return CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

func toCategoryResponses(categories []Category) []CategoryResponse {
	resp := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		resp[i] = toCategoryResponse(c)
	}
	return resp
}

func toCategoryFromRequest(req CategoryRequest) Category {
	return Category{
		Name: req.Name,
	}
}
