package routes

import (
	"net/http"
	"net/http/httptest"
	"shop-api/generated/db"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"
)

func TestGetUserId(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("test"), nil)
	token, err := GenerateToken(1, db.RoleMember, tokenAuth)
	if err != nil {
		t.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)
	r.Get("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, int64(1), GetUserId(w, r))
	})
	req := httptest.NewRequest("GET", "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token.Token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
}

func TestGetUserRole(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("test"), nil)
	token, err := GenerateToken(1, db.RoleMember, tokenAuth)
	if err != nil {
		t.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)
	r.Get("/api/users/me", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, db.RoleMember, GetUserRole(w, r))
	})
	req := httptest.NewRequest("GET", "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token.Token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
}
