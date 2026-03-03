package main

import (
	"shop-api/internal/api"
	_ "shop-api/internal/docs"
)

// @title Shop Api
// @version 1.0
// @description Swagger for Shop Api
// @BasePath /api
func main() {
	api.Run()
}
