package routes

import (
	"errors"
	"net/http"

	"books_api/repository"
	"books_api/service"

	// Adicione a implementação de AuthService se ela estiver ausente no pacote service
	// Se já existir, verifique se está definida corretamente neste local
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthRoutes configura as rotas de autenticação da aplicação.
// SetupAuthRoutes configura as rotas de autenticação da aplicação.
func SetupAuthRoutes(router *gin.Engine, authService *service.AuthService) {
	configureAuthGroup(router, authService)
}

// configureAuthGroup organiza as rotas dentro do grupo "/auth".
func configureAuthGroup(router *gin.Engine, authService *service.AuthService) {
	group := router.Group("/auth")
	{
		group.POST("/login", loginHandler(authService))
		group.POST("/register", registerHandler(authService))
	}
}

// loginHandler realiza a autenticação do usuário.
func loginHandler(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Requisição inválida"})
			return
		}

		token, err := authService.Authenticate(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// registerHandler realiza o registro de um novo usuário.
func registerHandler(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Requisição inválida"})
			return
		}

		if err := authService.Register(req.Username, req.Password); err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro interno no servidor"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
	}
}
