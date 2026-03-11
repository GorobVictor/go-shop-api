package main

import (
	"context"
	"log"
	"math/rand"
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

	// userRepo := repositories.NewUserRepository(conn)
	productRepo := repositories.NewProductRepository(conn)

	// for i := 0; i < 1000; i++ {
	// 	userRepo.CreateUser(context.Background(), db.CreateUserParams{
	// 		FirstName:    gofakeit.FirstName(),
	// 		LastName:     gofakeit.LastName(),
	// 		Email:        gofakeit.Email(),
	// 		PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
	// 		UserRole:     db.RoleMember,
	// 	})
	// }

	// fakeImage := pgtype.Text{String: "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg", Valid: true}

	// for i := 0; i < 100000; i++ {
	// 	prod := gofakeit.Product()
	// 	productRepo.CreateProduct(context.Background(), db.CreateProductParams{
	// 		Name:        prod.Name,
	// 		Description: pgtype.Text{String: prod.Description, Valid: true},
	// 		Price:       int64(prod.Price) * 100,
	// 		Discount:    0,
	// 		Image:       fakeImage,
	// 	})
	// }

	products, err := productRepo.GetProducts(context.Background(), 100, 0)
	receiptRepo := repositories.NewReceiptRepository(conn)
	if err != nil {
		log.Fatalln("problem with getting products, " + err.Error())
	}

	for i := 0; i < 20; i++ {
		min := randRange(0, 49)
		max := randRange(51, 99)
		randomProducts := products[min:max]
		sumPrice := int64(0)
		sumDiscount := int64(0)
		productsParams := make([]db.CreateReceiptProductParams, len(randomProducts))
		for index, product := range randomProducts {
			q := int32(gofakeit.Number(1, 10))
			sumPrice += product.Price * int64(q)
			sumDiscount += product.Discount * int64(q)
			productsParams[index] = db.CreateReceiptProductParams{
				ReceiptID: 0,
				ProductID: product.ID,
				Quantity:  q,
				Price:     product.Price,
				Discount:  product.Discount,
			}
		}

		_, _, err = receiptRepo.CreateReceipt(context.Background(), db.CreateReceiptParams{
			UserID:       2,
			SumPrice:     sumPrice,
			SumDiscount:  sumDiscount,
			StripeID:     pgtype.Text{String: gofakeit.UUID(), Valid: true},
			StripeStatus: db.StripeStatusSucceeded,
		}, productsParams)
		if err != nil {
			log.Fatalln("problem with creating receipt, " + err.Error())
		}
	}
}

func randRange(min, max int) int {
	return min + rand.Intn(max-min)
}
