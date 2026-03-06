package main

import (
	"shop-api/internal/api"
	_ "shop-api/internal/docs"
)

// @title Shop Api
// @version 1.0
// @description Swagger for Shop Api
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your token.
func main() {
	api.Run()
}
