package user

import (
	"errors"
	"log"
	"shop-api/generated/db"
	"shop-api/internal/database"
	"time"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(model SignInDto, tokenAuth *jwtauth.JWTAuth) (TokenDto, error) {
	ctx, conn, err := database.GetConnection()

	if err != nil {
		return TokenDto{}, err
	}

	defer conn.Close()

	q := db.New(conn)

	user, err := q.GetUserByEmail(ctx, model.Email)

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "no rows in result set" {
			return TokenDto{}, errors.New("incorrect email")
		}
	}

	if !checkPasswordHash(model.Password, user.PasswordHash) {
		return TokenDto{}, errors.New("incorrect password")
	}

	if err != nil {
		return TokenDto{}, err
	}

	token, err := generateToken(user, tokenAuth)

	return TokenDto{token}, err
}

func SignUp(model SignUpDto, tokenAuth *jwtauth.JWTAuth) (TokenDto, error) {
	ctx, conn, err := database.GetConnection()

	if err != nil {
		return TokenDto{}, err
	}

	defer conn.Close()

	q := db.New(conn)

	passHash, err := hashPassword(model.Password)

	if err != nil {
		return TokenDto{}, err
	}

	email, err := q.AnyEmail(ctx, model.Email)

	log.Println(email)

	if err != nil {
		log.Println(err.Error())
		if err.Error() != "no rows in result set" {
			return TokenDto{}, errors.New("Email already exists")
		}
	}

	user, err := q.CreateUser(ctx, db.CreateUserParams{
		FirstName: model.FirstName, LastName: model.LastName, Email: model.Email, PasswordHash: passHash, UserRole: db.RoleMember,
	})

	if err != nil {
		return TokenDto{}, err
	}

	token, err := generateToken(user, tokenAuth)

	return TokenDto{token}, err
}

func GetProfile(id int64) (ProfileDto, error) {
	ctx, conn, err := database.GetConnection()

	if err != nil {
		return ProfileDto{}, err
	}

	defer conn.Close()

	q := db.New(conn)

	user, err := q.GetUserProfile(ctx, id)

	return ProfileDto{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, UserRole: user.UserRole, CreatedAt: user.CreatedAt.Time}, err
}

func GetUsers(limit int64, offset int64) (result []ProfileDto, err error) {
	ctx, conn, err := database.GetConnection()

	if err != nil {
		return []ProfileDto{}, err
	}

	defer conn.Close()

	q := db.New(conn)

	users, err := q.GetUsers(ctx, db.GetUsersParams{Limit: limit, Offset: offset})

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
