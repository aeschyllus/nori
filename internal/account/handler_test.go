package account

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
	listFn   func(ctx context.Context) ([]Account, error)
	getFn    func(ctx context.Context, id int64) (Account, error)
	createFn func(ctx context.Context, a Account) (Account, error)
	updateFn func(ctx context.Context, id int64, a Account) (Account, error)
}

func (s *stubService) List(ctx context.Context) ([]Account, error) {
	return s.listFn(ctx)
}

func (s *stubService) GetByID(ctx context.Context, id int64) (Account, error) {
	return s.getFn(ctx, id)
}

func (s *stubService) Create(ctx context.Context, a Account) (Account, error) {
	return s.createFn(ctx, a)
}

func (s *stubService) Update(ctx context.Context, id int64, a Account) (Account, error) {
	return s.updateFn(ctx, id, a)
}

func setupRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAccountHandler(svc)
	g := r.Group("/accounts")
	g.GET("/", h.ListAccounts)
	g.GET("/:id", h.GetAccount)
	g.POST("/", h.CreateAccount)
	g.PUT("/:id", h.UpdateAccount)
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

func TestListAccountsOK(t *testing.T) {
	svc := &stubService{
		listFn: func(ctx context.Context) ([]Account, error) {
			return []Account{{ID: 1, Name: "Savings", Amount: 100}}, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/accounts/", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp []AccountResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 1 || resp[0].ID != 1 || resp[0].Amount != 100 {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestListAccountsEmptyIsArray(t *testing.T) {
	svc := &stubService{
		listFn: func(ctx context.Context) ([]Account, error) {
			return nil, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/accounts/", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Body.String(); got != "[]" {
		t.Fatalf("body = %q, want %q", got, "[]")
	}
}

func TestGetAccountRejectsBadID(t *testing.T) {
	svc := &stubService{
		getFn: func(ctx context.Context, id int64) (Account, error) {
			return Account{}, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/accounts/abc", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGetAccountNotFound(t *testing.T) {
	svc := &stubService{
		getFn: func(ctx context.Context, id int64) (Account, error) {
			return Account{}, ErrAccountNotFound
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodGet, "/accounts/999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestCreateAccountOK(t *testing.T) {
	svc := &stubService{
		createFn: func(ctx context.Context, a Account) (Account, error) {
			a.ID = 42
			return a, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPost, "/accounts/", map[string]any{
		"name": "New", "amountInCents": 1,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"accountId", "name", "amountInCents"} {
		if _, ok := resp[key]; !ok {
			t.Fatalf("response missing key %q: %s", key, w.Body.String())
		}
	}
}

func TestCreateAccountInvalidBody(t *testing.T) {
	svc := &stubService{
		createFn: func(ctx context.Context, a Account) (Account, error) {
			t.Error("createFn should not be called when binding fails")
			return Account{}, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPost, "/accounts/", map[string]any{
		"name": "", "amountInCents": 0,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestCreateAccountServiceRejects(t *testing.T) {
	svc := &stubService{
		createFn: func(ctx context.Context, a Account) (Account, error) {
			return Account{}, ErrInvalidAccount
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPost, "/accounts/", map[string]any{
		"name": "Valid", "amountInCents": 0,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateAccountOK(t *testing.T) {
	svc := &stubService{
		updateFn: func(ctx context.Context, id int64, a Account) (Account, error) {
			a.ID = id
			return a, nil
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPut, "/accounts/1", map[string]any{
		"name": "Renamed", "amountInCents": 999,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"accountId", "name", "amountInCents"} {
		if _, ok := resp[key]; !ok {
			t.Fatalf("response missing key %q: %s", key, w.Body.String())
		}
	}
}

func TestUpdateAccountNotFound(t *testing.T) {
	svc := &stubService{
		updateFn: func(ctx context.Context, id int64, a Account) (Account, error) {
			return Account{}, ErrAccountNotFound
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPut, "/accounts/999", map[string]any{
		"name": "X", "amountInCents": 1,
	})
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateAccountServiceRejects(t *testing.T) {
	svc := &stubService{
		updateFn: func(ctx context.Context, id int64, a Account) (Account, error) {
			return Account{}, ErrInvalidAccount
		},
	}
	r := setupRouter(svc)

	w := doJSON(t, r, http.MethodPut, "/accounts/999", map[string]any{
		"name": "X", "amountInCents": 1,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}
