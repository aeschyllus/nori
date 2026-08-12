package category

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubService struct {
	listFn   func(ctx context.Context) ([]Category, error)
	getFn    func(ctx context.Context, id int64) (Category, error)
	createFn func(ctx context.Context, c Category) (Category, error)
	updateFn func(ctx context.Context, id int64, c Category) (Category, error)
}

func (s *stubService) List(ctx context.Context) ([]Category, error) {
	return s.listFn(ctx)
}

func (s *stubService) GetByID(ctx context.Context, id int64) (Category, error) {
	return s.getFn(ctx, id)
}

func (s *stubService) Create(ctx context.Context, c Category) (Category, error) {
	return s.createFn(ctx, c)
}

func (s *stubService) Update(ctx context.Context, id int64, c Category) (Category, error) {
	return s.updateFn(ctx, id, c)
}

func setupRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCategoryHandler(svc)
	g := r.Group("/categories")
	g.GET("/", h.ListCategories)
	g.GET("/:id", h.GetCategory)
	g.POST("/", h.CreateCategory)
	g.PUT("/:id", h.UpdateCategory)
	return r
}

func doJSON(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestListCategoriesOK(t *testing.T) {
	svc := &stubService{
		listFn: func(ctx context.Context) ([]Category, error) {
			return []Category{{ID: 1, Name: "Food"}}, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/categories/", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp []CategoryResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 1 || resp[0].ID != 1 || resp[0].Name != "Food" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestListCategoryEmptyIsArray(t *testing.T) {
	svc := &stubService{
		listFn: func(ctx context.Context) ([]Category, error) {
			return nil, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/categories/", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Body.String(); got != "[]" {
		t.Fatalf("body = %q, want %q", got, "[]")
	}
}

func TestGetCategoryRejectsBadID(t *testing.T) {
	svc := &stubService{
		getFn: func(ctx context.Context, id int64) (Category, error) {
			return Category{}, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/categories/abc", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGetCategoryNotFound(t *testing.T) {
	svc := &stubService{
		getFn: func(ctx context.Context, id int64) (Category, error) {
			return Category{}, ErrCategoryNotFound
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/categories/999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestCreateCategoryOK(t *testing.T) {
	svc := &stubService{
		createFn: func(ctx context.Context, c Category) (Category, error) {
			c.ID = 42
			return c, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPost, "/categories/", map[string]any{
		"name": "Food",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"categoryId", "name"} {
		if _, ok := resp[key]; !ok {
			t.Fatalf("response missing key %q: %s", key, w.Body.String())
		}
	}
}

func TestCreateCategoryInvalidBody(t *testing.T) {
	svc := &stubService{
		createFn: func(ctx context.Context, c Category) (Category, error) {
			t.Error("createFn should not be called when binding fails")
			return Category{}, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPost, "/categories/", map[string]any{
		"name": "",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestCreateCategoryServiceRejects(t *testing.T) {
	svc := &stubService{
		createFn: func(ctx context.Context, c Category) (Category, error) {
			return Category{}, ErrInvalidCategory
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPost, "/categories/", map[string]any{
		"name": "Test",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateCategoryOK(t *testing.T) {
	svc := &stubService{
		updateFn: func(ctx context.Context, id int64, c Category) (Category, error) {
			c.ID = id
			return c, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPut, "/categories/1", map[string]any{
		"name": "Renamed",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"categoryId", "name"} {
		if _, ok := resp[key]; !ok {
			t.Fatalf("response key missing %q: %s", key, w.Body.String())
		}
	}
}

func TestUpdateCategoryNotFound(t *testing.T) {
	svc := &stubService{
		updateFn: func(ctx context.Context, id int64, c Category) (Category, error) {
			return Category{}, ErrCategoryNotFound
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPut, "/categories/999", map[string]any{
		"name": "X",
	})
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateCategoryServiceRejects(t *testing.T) {
	svc := &stubService{
		updateFn: func(ctx context.Context, id int64, c Category) (Category, error) {
			return Category{}, ErrInvalidCategory
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPut, "/categories/999", map[string]any{
		"name": "X",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}
