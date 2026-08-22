package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	pkgErrors "microservice/pkg/errors"
)

func newTestRouter(t *testing.T) *Router {
	t.Helper()
	mux := new(Router)
	mux.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if err, ok := recovered.(pkgErrors.StatusError); ok {
						w.WriteHeader(int(err.StatusCode()))
						fmt.Fprintf(w, "%s", err.Message())
						return
					}
					panic(recovered)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})
	return mux
}

func TestRouterBasicRoute(t *testing.T) {
	mux := newTestRouter(t)
	mux.Handle("/hello/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), "GET")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/hello", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRouterNotFound(t *testing.T) {
	mux := newTestRouter(t)
	mux.Handle("/hello/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), "GET")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/does-not-exist", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	mux := newTestRouter(t)
	mux.Handle("/hello/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), "GET")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/hello", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestRouterMultipleMethodsOnSamePattern(t *testing.T) {
	mux := newTestRouter(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("/multi/", handler, "GET")
	mux.Handle("/multi/", handler, "POST")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/multi", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRouterIntVariableRecorded(t *testing.T) {
	type ctxKey string

	var got int64
	mux := newTestRouter(t)
	mux.Handle("/users/{id:int}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Context().Value(ReturnContextKey("id")).(int64)
		w.WriteHeader(http.StatusOK)
	}), "GET")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/42", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got != 42 {
		t.Fatalf("expected id=42 in context, got %d", got)
	}
}

func TestRouterStringVariableRecorded(t *testing.T) {
	var got string
	mux := newTestRouter(t)
	mux.Handle("/files/{name}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Context().Value(ReturnContextKey("name")).(string)
		w.WriteHeader(http.StatusOK)
	}), "GET")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/files/report.pdf", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got != "report.pdf" {
		t.Fatalf("expected name='report.pdf' in context, got %q", got)
	}
}

func TestRouterTrailingSlashNormalized(t *testing.T) {
	mux := newTestRouter(t)
	called := false
	mux.Handle("/no-slash/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}), "GET")

	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/no-slash/", nil))
	if !called {
		t.Fatal("expected route to be called with normalized trailing slash")
	}
}

func TestRouterInvalidMethodPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on invalid method")
		}
	}()
	new(Router).Handle("/x/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), "FETCH")
}

func TestRouterDuplicateRegistrationPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate registration")
		}
	}()
	mux := new(Router)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	mux.Handle("/dup/", h, "GET")
	mux.Handle("/dup/", h, "GET")
}

func TestConvertToType(t *testing.T) {
	cases := map[string]variableType{
		"int":   variableInt,
		"float": variableFloat,
		"other": variableString,
	}
	for in, want := range cases {
		if got := convertToType(in); got != want {
			t.Errorf("convertToType(%q) = %v, want %v", in, got, want)
		}
	}
}
