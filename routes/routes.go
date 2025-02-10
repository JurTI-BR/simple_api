package routes

import (
	"books_api/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(router *gin.Engine, authService *service.AuthService) {
	// Configura as rotas de autenticação

	AuthRoutes(router, authService)

	// Configura as rotas de livros
	BookRoutes(router)
}
