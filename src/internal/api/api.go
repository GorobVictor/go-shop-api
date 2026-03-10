package api

import (
	"log"
	"net/http"
	"os"
	"shop-api/internal/api/routes"
	"shop-api/internal/database"
	"shop-api/internal/database/repositories"
	"shop-api/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Run() {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Fatalln("JWT secret is null")
	}
	conn, err := database.GetConnection()

	if err != nil {
		log.Fatalln("problem with db pool, " + err.Error())
	}

	defer conn.Close()

	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	userRepo := repositories.NewUserRepository(conn)
	userSvc := user.NewUserService(userRepo)
	userHandler := routes.NewUserHandler(userSvc, tokenAuth)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	initSwagger(r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	})

	userHandler.Users(r)

	log.Println("Server starting on port :3000")
	log.Fatalln(http.ListenAndServe(":3000", r))
}

func initSwagger(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))
}
