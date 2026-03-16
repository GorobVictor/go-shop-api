package routes

import (
	"encoding/json"
	"net/http"
	"shop-api/generated/db"
	customerrors "shop-api/internal/custom_errors"
	"strconv"
	"time"

	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
)

func GenerateToken(id int64, role db.Role, tokenAuth *jwtauth.JWTAuth) (response TokenDto, err error) {
	_, response.Token, err = tokenAuth.Encode(map[string]interface{}{
		"user_id":   id,
		"user_role": role,
	})

	return response, err
}

// Token model
type TokenDto struct {
	Token string `json:"token"`
}

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

func GetCORSMiddleware(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})(next)
}

func GetPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch v := err.(type) {
				case customerrors.BadRequestError:
					WriteError(w, http.StatusBadRequest, &v)
				case customerrors.UnauthorizedError:
					WriteError(w, http.StatusUnauthorized, &v)
				case customerrors.ForbiddenError:
					WriteError(w, http.StatusForbidden, &v)
				case customerrors.InternalServerError:
					WriteError(w, http.StatusInternalServerError, &v)
				case *customerrors.InternalServerError:
					WriteError(w, http.StatusInternalServerError, v)
				default:
					WriteError(w, http.StatusInternalServerError, customerrors.NewInternalServerError())
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

func ReadBody(w http.ResponseWriter, r *http.Request, v any) {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
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

func WriteOkResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func WriteInternalServerError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		json.NewEncoder(w).Encode(customerrors.NewInternalServerError())
	}
}
