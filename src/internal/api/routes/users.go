package routes

import (
	"encoding/json"
	"net/http"
	"shop-api/internal/api/middlewares"
	"shop-api/internal/usecase/user"

	"github.com/go-chi/chi/v5"
)

func Users(r *chi.Mux, clerkKey string) {
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/signin", signIn)
		r.Post("/signup", signUp)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.GetAuthMiddleware)
			r.Get("/me", me)
		})
	})
}

// Sign In
// @Summary Sign in
// @Tags users
// @Param user body user.SignInDto true "User details for login"
// @Success 200 {object} user.TokenDto
// @Router /users/signin [post]
func signIn(w http.ResponseWriter, r *http.Request) {

	var model user.SignInDto
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	user, err := user.SignIn(model)

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
func signUp(w http.ResponseWriter, r *http.Request) {
	var model user.SignUpDto
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	user, err := user.SignUp(model)

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(user)
}

// Get my profile
// @Summary Get my profile
// @Tags users
// @Success 200 {object} user.ProfileDto
// @Router /users/me [get]
func me(w http.ResponseWriter, r *http.Request) {

	// log.Println(claims)

	user, err := user.GetProfile()

	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	json.NewEncoder(w).Encode(user)
}
