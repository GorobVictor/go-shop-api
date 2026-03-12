package routes

import (
	"encoding/json"
	"net/http"
	"shop-api/generated/db"
	customerrors "shop-api/internal/custom_errors"
	"strconv"
	"time"

	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
)

func GetUserId(w http.ResponseWriter, r *http.Request) int64 {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userIdValue, ok := claims["user_id"]

	if !ok {
		panic(customerrors.UnauthorizedError{Message: "user_id not found"})
	}

	userId, ok := userIdValue.(float64)

	if !ok {
		panic(customerrors.UnauthorizedError{Message: "user_id not found"})
	}

	return int64(userId)
}

func GetUserRole(w http.ResponseWriter, r *http.Request) db.Role {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userRoleValue, ok := claims["user_role"]

	if !ok {
		panic(customerrors.ForbiddenError{Message: "user_role not found"})
	}

	userRole, ok := userRoleValue.(string)

	if !ok {
		panic(customerrors.ForbiddenError{Message: "user_role not found"})
	}

	return db.Role(userRole)
}

func GetAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetUserRole(w, r)
		if role != db.RoleAdmin {
			panic(customerrors.ForbiddenError{Message: "You are not admin!"})
		}
		next.ServeHTTP(w, r)
	})
}

func GetRateLimitMiddleware(next http.Handler) http.Handler {
	return httprate.LimitByIP(100, time.Minute)(next)
}

func GetPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch v := err.(type) {
				case customerrors.BadRequestError:
					json.NewEncoder(w).Encode(v)
					w.WriteHeader(http.StatusBadRequest)
				case customerrors.UnauthorizedError:
					json.NewEncoder(w).Encode(v)
					w.WriteHeader(http.StatusUnauthorized)
				case customerrors.ForbiddenError:
					json.NewEncoder(w).Encode(v)
					w.WriteHeader(http.StatusForbidden)
				case customerrors.InternalServerError:
					json.NewEncoder(w).Encode(v)
					w.WriteHeader(http.StatusInternalServerError)
				default:
					json.NewEncoder(w).Encode(customerrors.NewInternalServerError())
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ReadBody(w http.ResponseWriter, r *http.Request, v any) {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		panic(customerrors.BadRequestError{Message: err.Error()})
	}
}

func GetQueryInt32(r *http.Request, name string) int32 {
	queries := r.URL.Query()
	value, err := strconv.ParseInt(queries.Get(name), 10, 32)
	if err != nil {
		panic(customerrors.BadRequestError{Message: err.Error()})
	}
	return int32(value)
}
