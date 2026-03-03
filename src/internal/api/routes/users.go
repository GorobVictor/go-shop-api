package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Users(r *chi.Mux) {
	r.Get("/api/users/{userID}", getById)
}

// GetUser godoc
// @Summary Get user
// @Tags users
// @Param id path int true "User ID"
// @Router /users/{id} [get]
func getById(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	w.Write([]byte("Fetching user: " + userID))
}
