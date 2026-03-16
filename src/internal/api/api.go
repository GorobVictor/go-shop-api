package api

import (
	"context"
	"log"
	"net/http"
	"shop-api/generated/db"
	"shop-api/internal/api/routes"
	"shop-api/internal/config"
	customerrors "shop-api/internal/custom_errors"
	"shop-api/internal/database/repositories"
	"shop-api/internal/usecase/product"
	"shop-api/internal/usecase/receipt"
	"shop-api/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v84"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("problem with config, " + err.Error())
	}

	conn, err := pgxpool.New(context.Background(), config.PostgresUrl)
	if err != nil {
		log.Fatalln("problem with db pool, " + err.Error())
	}
	defer conn.Close()

	q := db.New(conn)

	tokenAuth := jwtauth.New("HS256", []byte(config.JwtSecret), nil)
	stripeClient := stripe.NewClient(config.StripeSecret)

	userRepo := repositories.NewUserRepository(conn, q)
	productRepo := repositories.NewProductRepository(conn, q)
	receiptRepo := repositories.NewReceiptRepository(conn, q)

	userSvc := user.NewUserService(userRepo)
	productSvc := product.NewProductService(productRepo)
	receiptSvc := receipt.NewReceiptService(receiptRepo, productRepo, stripeClient, &config)

	userHandler := routes.NewUserHandler(userSvc, tokenAuth)
	productHandler := routes.NewProductHandler(productSvc, tokenAuth)
	paymentHandler := routes.NewPaymentHandler(receiptSvc, tokenAuth)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(routes.GetPanicMiddleware)
	r.Use(routes.GetRateLimitMiddleware)
	r.Use(routes.GetCORSMiddleware)

	initSwagger(r)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		err := conn.Ping(context.Background())
		if err != nil {
			panic(customerrors.NewInternalServerError())
		}
		w.Write([]byte("Welcome!"))
	})

	userHandler.Users(r)
	productHandler.Products(r)
	paymentHandler.Payment(r)

	log.Println("Server starting on port :" + config.Port)
	log.Fatalln(http.ListenAndServe(":"+config.Port, r))
}

func initSwagger(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))
}
