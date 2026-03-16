package repositories

import (
	"context"
	"errors"
	"log"
	"shop-api/generated/db"
	customerrors "shop-api/internal/custom_errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
	q  *db.Queries
}

func NewUserRepository(db *pgxpool.Pool, q *db.Queries) *UserRepository {
	return &UserRepository{db: db, q: q}
}

func (r *UserRepository) GetUserProfile(ctx context.Context, id int64) (db.GetUserProfileRow, error) {
	user, err := r.q.GetUserProfile(ctx, id)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			panic(&customerrors.BadRequestError{Message: "wrong user id"})
		}
	}

	return user, err
}

func (r *UserRepository) GetUsers(ctx context.Context, limit int32, offset int32) ([]db.GetUsersRow, error) {
	users, err := r.q.GetUsers(ctx, db.GetUsersParams{Limit: limit, Offset: offset})

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return users, nil
		}
	}

	return users, err
}

func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)

	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			panic(&customerrors.BadRequestError{Message: "incorrect email"})
		}
	}

	return user, err
}

func (r *UserRepository) AnyEmail(ctx context.Context, email string) error {
	_, err := r.q.AnyEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
	}

	if err == nil {
		panic(&customerrors.BadRequestError{Message: "Email already exists"})
	}

	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, user db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, user)
}
