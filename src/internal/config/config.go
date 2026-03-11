package config

import "github.com/caarlos0/env/v11"

type Config struct {
	PostgresUrl  string `env:"POSTGRES_CONN_STR"`
	JwtSecret    string `env:"JWT_SECRET"`
	StripeSecret string `env:"STRIPE_SECRET"`
	FrontendUrl  string `env:"FRONTEND_URL"`
	BackendUrl   string `env:"BACKEND_URL"`
}

func LoadConfig() (Config, error) {
	return env.ParseAs[Config]()
}
