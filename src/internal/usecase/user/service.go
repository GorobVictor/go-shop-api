package user

import (
	"context"
	"shop-api/generated/db"
	customerrors "shop-api/internal/custom_errors"
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

func (s *UserService) SignIn(ctx context.Context, model SignInDto, tokenAuth *jwtauth.JWTAuth) (SignInResDto, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, model.Email)

	if err != nil {
		return SignInResDto{}, err
	}

	if !checkPasswordHash(model.Password, user.PasswordHash) {
		panic(customerrors.BadRequestError{Message: "incorrect password"})
	}

	return NewSignInResDto(user), err
}

func (s *UserService) SignUp(ctx context.Context, model SignUpDto, tokenAuth *jwtauth.JWTAuth) (SignInResDto, error) {
	passHash, err := hashPassword(model.Password)

	if err != nil {
		return SignInResDto{}, err
	}

	err = s.userRepo.AnyEmail(ctx, model.Email)

	if err != nil {
		return SignInResDto{}, err
	}

	user, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
		FirstName: model.FirstName, LastName: model.LastName, Email: model.Email, PasswordHash: passHash, UserRole: db.RoleMember,
	})

	if err != nil {
		return SignInResDto{}, err
	}

	return NewSignInResDto(user), err
}

func (s *UserService) GetProfile(ctx context.Context, id int64) (ProfileDto, error) {

	user, err := s.userRepo.GetUserProfile(ctx, id)

	return ProfileDto{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, UserRole: user.UserRole, CreatedAt: user.CreatedAt.Time}, err
}

func (s *UserService) GetUsers(ctx context.Context, limit int32, offset int32) (result UsersPaginationDto, err error) {
	result.Limit = limit
	result.Offset = offset

	result.Total, err = s.userRepo.CountUsers(ctx)

	if err != nil {
		return result, err
	}

	users, err := s.userRepo.GetUsers(ctx, limit, offset)

	for _, user := range users {
		result.Users = append(result.Users, ProfileDto{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, UserRole: user.UserRole, CreatedAt: user.CreatedAt.Time})
	}

	return result, err
}

// SignIdDto represents the request body for sign in a user
type SignInDto struct {
	Email    string `json:"email" example:"test@test.com"`
	Password string `json:"password" example:"test123"`
}

type UsersPaginationDto struct {
	Users  []ProfileDto `json:"users"`
	Total  int64        `json:"total"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
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

type SignInResDto struct {
	ID       int64   `json:"id" example:"1"`
	UserRole db.Role `json:"userRole" example:"member"`
}

func NewSignInResDto(user db.User) SignInResDto {
	return SignInResDto{ID: user.ID, UserRole: user.UserRole}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
