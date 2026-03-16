package routes

import (
	"context"
	"net/http"
	"shop-api/internal/usecase/user"

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
	ReadBody(w, r, &model)

	user, err := h.userSvc.SignIn(context.Background(), model, h.tokenAuth)

	if err != nil {
		CheckError(w, err)
		return
	}

	response, err := GenerateToken(user.ID, user.UserRole, h.tokenAuth)

	if err != nil {
		CheckError(w, err)
		return
	}

	WriteOkResponse(w, response)
}

// Sign Up
// @Summary Sign Up
// @Tags users
// @Param user body user.SignUpDto true "User details for registration"
// @Success 200 {object} user.TokenDto
// @Router /users/signup [post]
func (h *UserHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var model user.SignUpDto
	ReadBody(w, r, &model)

	user, err := h.userSvc.SignUp(context.Background(), model, h.tokenAuth)

	if err != nil {
		CheckError(w, err)
		return
	}

	response, err := GenerateToken(user.ID, user.UserRole, h.tokenAuth)

	if err != nil {
		CheckError(w, err)
		return
	}

	WriteOkResponse(w, response)
}

// Get my profile
// @Summary Get my profile
// @Tags users
// @Security ApiKeyAuth
// @Success 200 {object} user.ProfileDto
// @Router /users/me [get]
func (h *UserHandler) me(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserId(w, r)
	if err != nil {
		CheckError(w, err)
		return
	}

	user, err := h.userSvc.GetProfile(context.Background(), userId)

	if err != nil {
		CheckError(w, err)
		return
	}

	WriteOkResponse(w, user)
}

// Get users
// @Summary Get my profile
// @Tags users
// @Security ApiKeyAuth
// @Success 200 {object} user.UsersPaginationDto
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Router /users/get [get]
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	result, err := h.userSvc.GetUsers(context.Background(), GetQueryInt32(r, "limit"), GetQueryInt32(r, "offset"))

	if err != nil {
		CheckError(w, err)
		return
	}

	WriteOkResponse(w, result)
}
