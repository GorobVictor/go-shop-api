package routes

import (
	"encoding/json"
	"net/http"
	"shop-api/generated/db"

	"github.com/go-chi/jwtauth"
)

func GetUserId(w http.ResponseWriter, r *http.Request) int64 {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userIdValue, ok := claims["user_id"]

	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
	}

	userId, ok := userIdValue.(float64)

	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
	}

	return int64(userId)
}

func GetUserRole(w http.ResponseWriter, r *http.Request) db.Role {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userRoleValue, ok := claims["user_role"]

	if !ok {
		http.Error(w, "user_role not found", http.StatusForbidden)
	}

	userRole, ok := userRoleValue.(string)

	if !ok {
		http.Error(w, "user_role not found", http.StatusForbidden)
	}

	return db.Role(userRole)
}

func GetAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(w, r)
		if role != db.RoleAdmin {
			writeForbiddenStr(w, "You are not admin!")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeBadRequest(w http.ResponseWriter, err error) {
	writeError(w, err, http.StatusBadRequest)
}

func writeForbidden(w http.ResponseWriter, err error) {
	writeError(w, err, http.StatusForbidden)
}

func writeForbiddenStr(w http.ResponseWriter, err string) {
	writeErrorStr(w, err, http.StatusForbidden)
}

func writeError(w http.ResponseWriter, err error, code int) {
	if err != nil {
		json.NewEncoder(w).Encode(ResError{err.Error()})
	}
	w.WriteHeader(code)
}

func writeErrorStr(w http.ResponseWriter, err string, code int) {
	json.NewEncoder(w).Encode(ResError{err})
	w.WriteHeader(code)
}

type ResError struct {
	Message string `json:"message" example:"Error"`
}
