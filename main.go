package main

import (
	"books_api/config"
	. "books_api/middleware"
	"books_api/repository"
	"books_api/routes"
	"books_api/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Não foi possível carregar o arquivo .env, utilizando variáveis de ambiente padrão.")
	}

	config.ConnectDatabase()
	config.ConnectRedis()

	r := gin.Default()

	r.Use(CORSMiddleware())

	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)

	routes.SetupRoutes(r, authService)

	port := ":8080"
	log.Println("Servidor rodando na porta", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
