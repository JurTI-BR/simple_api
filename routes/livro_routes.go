package routes

import (
	"books_api/middleware"
	"books_api/models"
	"books_api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine) {
	livroService := service.NewLivroService() // Criar uma instância única do serviço

	livros := router.Group("/livros")
	livros.Use(middleware.AuthMiddleware())
	{
		livros.GET("", func(c *gin.Context) { listarLivros(c, livroService) })
		livros.GET("/:id", func(c *gin.Context) { buscarLivroPorID(c, livroService) })
		livros.POST("", func(c *gin.Context) { criarLivro(c, livroService) })
		livros.PUT("/:id", func(c *gin.Context) { atualizarLivro(c, livroService) })
		livros.DELETE("/:id", func(c *gin.Context) { deletarLivro(c, livroService) })
	}
}

func getIDFromParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	return uint(id), err
}

func listarLivros(c *gin.Context, srv service.LivroService) {
	ctx := c.Request.Context()

	livros, err := srv.ListarLivros(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar livros"})
		return
	}

	c.JSON(http.StatusOK, livros)
}

// buscarLivroPorID busca um livro pelo seu ID.
func buscarLivroPorID(c *gin.Context, srv service.LivroService) {
	ctx := c.Request.Context()

	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	livro, err := srv.BuscarLivroPorID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar livro"})
		return
	}

	if livro == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Livro não encontrado"})
		return
	}

	c.JSON(http.StatusOK, livro)
}

func criarLivro(c *gin.Context, srv service.LivroService) {
	ctx := c.Request.Context()

	var novoLivro models.Livro
	if err := c.ShouldBindJSON(&novoLivro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos"})
		return
	}

	if err := srv.CriarLivro(ctx, &novoLivro); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar livro"})
		return
	}

	c.JSON(http.StatusCreated, novoLivro)
}

func atualizarLivro(c *gin.Context, srv service.LivroService) {
	ctx := c.Request.Context()

	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	var livroAtualizado models.Livro
	if err := c.ShouldBindJSON(&livroAtualizado); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos"})
		return
	}

	livro, err := srv.AtualizarLivro(ctx, id, &livroAtualizado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar livro"})
		return
	}
	if livro == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Livro não encontrado"})
		return
	}

	c.JSON(http.StatusOK, livro)
}

func deletarLivro(c *gin.Context, srv service.LivroService) {
	ctx := c.Request.Context()

	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	if err := srv.DeletarLivro(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao deletar livro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Livro deletado com sucesso"})
}
