package user

import (
	"errors"
	"log"
	"shop-api/generated/db"
	"shop-api/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func SignIn(model SignInDto) (TokenDto, error) {
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

	log.Println(user)

	return TokenDto{"test token"}, nil
}

func SignUp(model SignUpDto) (TokenDto, error) {
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
	if email != "" && err == nil {
		return TokenDto{}, errors.New("Email already exists")
	}

	if err != nil {
		return TokenDto{}, err
	}

	user, err := q.CreateUser(ctx, db.CreateUserParams{
		FirstName: model.FirstName, LastName: model.LastName, Email: model.Email, PasswordHash: passHash,
	})

	if err != nil {
		return TokenDto{}, err
	}

	log.Println(user)

	return TokenDto{"test token"}, nil
}

func GetProfile() (ProfileDto, error) {
	// ctx, conn, err := database.GetConnection()

	// if err != nil {
	// 	return ProfileDto{}, err
	// }

	// defer conn.Close()

	// q := db.New(conn)

	// email, err := q.AnyEmail(ctx, model.Email)

	// log.Println(email)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	if err.Error() != "no rows in result set" {
	// 		return TokenDto{}, errors.New("Email already exists")
	// 	}
	// }
	// if email != "" && err == nil {
	// 	return TokenDto{}, errors.New("Email already exists")
	// }

	// user, err := q.CreateUser(ctx, db.CreateUserParams{FirstName: model.FirstName, LastName: model.LastName, Email: model.Email, PasswordHash: passHash})

	// if err != nil {
	// 	return TokenDto{}, err
	// }

	// log.Println(user)

	return ProfileDto{}, nil
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
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Wick"`
	Email     string `json:"email" example:"test@test.com"`
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
