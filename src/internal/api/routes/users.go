package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"shop-api/internal/usecase/user"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	userSvc   *user.UserService
	tokenAuth *jwtauth.JWTAuth
}

func NewUserHandler(userSvc *user.UserService, tokenAuth *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{userSvc: userSvc, tokenAuth: tokenAuth}
}

func (h *UserHandler) Users(r *chi.Mux) {
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/signin", h.signIn)
		r.Post("/signup", h.signUp)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Get("/me", h.me)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(h.tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Use(GetAdminMiddleware)
			r.Get("/get", h.getUsers)
		})
	})
}

// Sign In
// @Summary Sign in
// @Tags users
// @Param user body user.SignInDto true "User details for login"
// @Success 200 {object} user.TokenDto
// @Router /users/signin [post]
func (h *UserHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var model user.SignInDto
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	user, err := h.userSvc.SignIn(context.Background(), model, h.tokenAuth)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(user)
}

// Sign Up
// @Summary Sign Up
// @Tags users
// @Param user body user.SignUpDto true "User details for registration"
// @Success 200 {object} user.TokenDto
// @Router /users/signup [post]
func (h *UserHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var model user.SignUpDto
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	user, err := h.userSvc.SignUp(context.Background(), model, h.tokenAuth)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(user)
}

// Get my profile
// @Summary Get my profile
// @Tags users
// @Security ApiKeyAuth
// @Success 200 {object} user.ProfileDto
// @Router /users/me [get]
func (h *UserHandler) me(w http.ResponseWriter, r *http.Request) {

	userId := GetUserId(w, r)

	user, err := h.userSvc.GetProfile(context.Background(), userId)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(user)
}

// Get users
// @Summary Get my profile
// @Tags users
// @Security ApiKeyAuth
// @Success 200 {object} []user.ProfileDto
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Router /users/get [get]
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	limit, err := strconv.ParseInt(queries.Get("limit"), 10, 32)
	if err != nil {
		writeBadRequest(w, err)
	}
	offset, err := strconv.ParseInt(queries.Get("offset"), 10, 64)
	if err != nil {
		writeBadRequest(w, err)
	}

	result, err := h.userSvc.GetUsers(context.Background(), int32(limit), int32(offset))

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(result)
}
