package main

import (
	"context"
	"log"
	"shop-api/generated/db"
	"shop-api/internal/config"
	"shop-api/internal/database/repositories"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("problem with config, " + err.Error())
	}

	conn, err := pgxpool.New(context.Background(), config.PostgresUrl)
	if err != nil {
		log.Fatalln("problem with db pool, " + err.Error())
	}
	defer conn.Close()

	userRepo := repositories.NewUserRepository(conn)
	productRepo := repositories.NewProductRepository(conn)

	for i := 0; i < 1000; i++ {
		userRepo.CreateUser(context.Background(), db.CreateUserParams{
			FirstName:    gofakeit.FirstName(),
			LastName:     gofakeit.LastName(),
			Email:        gofakeit.Email(),
			PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
			UserRole:     db.RoleMember,
		})
	}

	fakeImage := pgtype.Text{String: "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg", Valid: true}

	for i := 0; i < 100000; i++ {
		prod := gofakeit.Product()
		productRepo.CreateProduct(context.Background(), db.CreateProductParams{
			Name:        prod.Name,
			Description: pgtype.Text{String: prod.Description, Valid: true},
			Price:       int64(prod.Price) * 100,
			Discount:    0,
			Image:       fakeImage,
		})
	}
}
