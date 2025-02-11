package main

import (
	"books_api/config"
	"books_api/repository"
	"books_api/routes"
	"books_api/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Não foi possível carregar o arquivo .env, utilizando variáveis de ambiente padrão.")
	}

	config.ConnectDatabase()
	config.ConnectRedis()
	// Criar instância do LivroService usando o banco PostgreSQL
	livroService := service.NewLivroService(config.DB)

	// Criar instância do UserService e AuthService
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)

	// Criar router do Gin
	r := gin.Default()

	// Permitir CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Servir arquivos de imagem
	r.Static("/uploads", "./uploads")

	// Configurar rotas passando `authService` e `livroService`
	routes.SetupRoutes(r, authService, livroService)

	// Iniciar servidor
	port := ":8080"
	log.Println("Servidor rodando na porta", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
