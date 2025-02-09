package main

import (
	"books_api/config"
	"books_api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Carregar variáveis de ambiente
	config.ConnectDatabase()
	config.ConnectRedis()

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
