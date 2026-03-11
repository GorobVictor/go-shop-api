package repositories

import (
	"context"
	"errors"
	"log"
	"shop-api/generated/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserProfile(ctx context.Context, id int64) (db.GetUserProfileRow, error) {
	q := db.New(r.db)
	user, err := q.GetUserProfile(ctx, id)

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "no rows in result set" {
			return user, errors.New("wrong user id")
		}
	}

	return user, err
}

func (r *UserRepository) GetUsers(ctx context.Context, limit int32, offset int32) ([]db.GetUsersRow, error) {
	q := db.New(r.db)
	users, err := q.GetUsers(ctx, db.GetUsersParams{Limit: limit, Offset: offset})

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "no rows in result set" {
			return users, nil
		}
	}

	return users, err
}

func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	q := db.New(r.db)
	return q.CountUsers(ctx)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	q := db.New(r.db)
	user, err := q.GetUserByEmail(ctx, email)

	if err != nil {
		log.Println(err.Error())
		if err.Error() == "no rows in result set" {
			return user, errors.New("incorrect email")
		}
	}

	return user, err
}

func (r *UserRepository) AnyEmail(ctx context.Context, email string) error {
	q := db.New(r.db)
	email, err := q.AnyEmail(ctx, email)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil
		}
	}

	if err == nil {
		return errors.New("Email already exists")
	}

	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, user db.CreateUserParams) (db.User, error) {
	q := db.New(r.db)

	return q.CreateUser(ctx, user)
}
