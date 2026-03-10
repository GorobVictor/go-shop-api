package user

import (
	"context"
	"errors"
	"shop-api/generated/db"
	"shop-api/internal/database/repositories"
	"time"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo}
}

func (s *UserService) SignIn(ctx context.Context, model SignInDto, tokenAuth *jwtauth.JWTAuth) (TokenDto, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, model.Email)

	if err != nil {
		return TokenDto{}, err
	}

	if !checkPasswordHash(model.Password, user.PasswordHash) {
		return TokenDto{}, errors.New("incorrect password")
	}

	token, err := generateToken(user, tokenAuth)

	return TokenDto{token}, err
}

func (s *UserService) SignUp(ctx context.Context, model SignUpDto, tokenAuth *jwtauth.JWTAuth) (TokenDto, error) {
	passHash, err := hashPassword(model.Password)

	if err != nil {
		return TokenDto{}, err
	}

	err = s.userRepo.AnyEmail(ctx, model.Email)

	if err != nil {
		return TokenDto{}, err
	}

	user, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
		FirstName: model.FirstName, LastName: model.LastName, Email: model.Email, PasswordHash: passHash, UserRole: db.RoleMember,
	})

	if err != nil {
		return TokenDto{}, err
	}

	token, err := generateToken(user, tokenAuth)

	return TokenDto{token}, err
}

func (s *UserService) GetProfile(ctx context.Context, id int64) (ProfileDto, error) {

	user, err := s.userRepo.GetUserProfile(ctx, id)

	return ProfileDto{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, UserRole: user.UserRole, CreatedAt: user.CreatedAt.Time}, err
}

func (s *UserService) GetUsers(ctx context.Context, limit int32, offset int32) (result []ProfileDto, err error) {

	users, err := s.userRepo.GetUsers(ctx, limit, offset)

	for _, user := range users {
		result = append(result, ProfileDto{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, UserRole: user.UserRole, CreatedAt: user.CreatedAt.Time})
	}

	return result, err
}

func generateToken(user db.User, tokenAuth *jwtauth.JWTAuth) (token string, err error) {
	_, token, err = tokenAuth.Encode(map[string]interface{}{
		"user_id":   user.ID,
		"user_role": user.UserRole,
	})

	return token, err
}

// SignIdDto represents the request body for sign in a user
type SignInDto struct {
	Email    string `json:"email" example:"test@test.com"`
	Password string `json:"password" example:"test123"`
}

// SignUpDto represents the request body for sign up a user
type SignUpDto struct {
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Wick"`
	Email     string `json:"email" example:"test@test.com"`
	Password  string `json:"password" example:"test123"`
}

type ProfileDto struct {
	ID        int64
	FirstName string    `json:"firstName" example:"John"`
	LastName  string    `json:"lastName" example:"Wick"`
	Email     string    `json:"email" example:"test@test.com"`
	UserRole  db.Role   `json:"userRole" example:"member"`
	CreatedAt time.Time `json:"createdAt"`
}

// Token model
type TokenDto struct {
	Token string `json:"token"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
