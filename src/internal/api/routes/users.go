package routes

import (
	"encoding/json"
	"net/http"
	"shop-api/internal/usecase/user"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

func Users(r *chi.Mux, tokenAuth *jwtauth.JWTAuth) {
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			signIn(w, r, tokenAuth)
		})
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			signUp(w, r, tokenAuth)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Get("/me", me)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Use(GetAdminMiddleware)
			r.Get("/get", getUsers)
		})
	})
}

// Sign In
// @Summary Sign in
// @Tags users
// @Param user body user.SignInDto true "User details for login"
// @Success 200 {object} user.TokenDto
// @Router /users/signin [post]
func signIn(w http.ResponseWriter, r *http.Request, tokenAuth *jwtauth.JWTAuth) {

	var model user.SignInDto
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	user, err := user.SignIn(model, tokenAuth)

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
func signUp(w http.ResponseWriter, r *http.Request, tokenAuth *jwtauth.JWTAuth) {
	var model user.SignUpDto
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	user, err := user.SignUp(model, tokenAuth)

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
func me(w http.ResponseWriter, r *http.Request) {

	userId := GetUserId(w, r)

	user, err := user.GetProfile(userId)

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
func getUsers(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	limit, err := strconv.ParseInt(queries.Get("limit"), 10, 64)
	if err != nil {
		writeBadRequest(w, err)
	}
	offset, err := strconv.ParseInt(queries.Get("offset"), 10, 64)
	if err != nil {
		writeBadRequest(w, err)
	}

	result, err := user.GetUsers(limit, offset)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(result)
}
