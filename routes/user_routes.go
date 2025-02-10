package routes

import (
	"books_api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRoutes configura as rotas de autenticação
func AuthRoutes(router *gin.Engine, authService *service.AuthService) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", func(c *gin.Context) {
			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
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
		})

		authGroup.POST("/register", func(c *gin.Context) {
			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Requisição inválida"})
				return
			}

			if err := authService.Register(req.Username, req.Password); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
		})
	}
}
