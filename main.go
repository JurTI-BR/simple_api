package main

import (
	"books_api/config"
	"books_api/repository"
	"books_api/routes"
	"books_api/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Carregar variáveis de ambiente
	config.ConnectDatabase()
	config.ConnectRedis()

	r := gin.Default()

	// Inicializa o repositório de usuários
	userRepository := repository.NewUserRepository(config.DB)

	// Inicializa o serviço de autenticação passando o repositório
	authService := service.NewAuthService(userRepository)

	// Configura todas as rotas da aplicação
	routes.SetupRoutes(r, authService)

	// Inicia o servidor
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
