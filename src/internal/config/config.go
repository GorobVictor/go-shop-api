package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	PostgresUrl  string `env:"POSTGRES_CONN_STR,required"`
	JwtSecret    string `env:"JWT_SECRET,required"`
	StripeSecret string `env:"STRIPE_SECRET,required"`
	FrontendUrl  string `env:"FRONTEND_URL" envDefault:"http://localhost:5173"`
	BackendUrl   string `env:"BACKEND_URL" envDefault:"http://localhost:3000"`
	Port         string `env:"PORT" envDefault:"3000"`
}

func LoadConfig() (config Config, err error) {
	return env.ParseAs[Config]()
}
