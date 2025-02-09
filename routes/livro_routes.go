package routes

import (
	"net/http"

	"books_api/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/livros", func(c *gin.Context) {
		livros, err := service.ListarLivros()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar livros"})
			return
		}
		c.JSON(http.StatusOK, livros)
	})
}
