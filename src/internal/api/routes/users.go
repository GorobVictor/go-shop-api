package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-api/internal/db"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, "user=shop_user password=shop_password dbname=shop_db sslmode=disable host=postgres port=5432")

	returnError(w, err)

	defer conn.Close()

	q := db.New(conn)

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 0, 64)

	returnError(w, err)

	user, err := q.GetUser(ctx, userID)

	returnError(w, err)

	json, err := json.Marshal(user)

	returnError(w, err)

	w.Write(json)
}

func returnError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
