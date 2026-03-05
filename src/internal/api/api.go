package api

import (
	"log"
	"net/http"
	"os"
	"shop-api/internal/api/routes"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Run() {
	clerkConnStr := os.Getenv("CLERK_KEY")

	if clerkConnStr == "" {
		log.Fatal("clerk connection string is empty")
	}

	clerk.SetKey(clerkConnStr)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	initSwagger(r)

	routes.Users(r, clerkConnStr)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	})

	log.Println("Server starting on port :3000")
	log.Fatalln(http.ListenAndServe(":3000", r))
}

func initSwagger(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))
}
