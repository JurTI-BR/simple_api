package routes

import (
	"books_api/middleware"
	"books_api/models"
	"books_api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookRoutes configura as rotas dos livros
func BookRoutes(router *gin.Engine) {
	livros := router.Group("/livros")
	livros.Use(middleware.AuthMiddleware())
	{
		livros.GET("", listarLivros)
		livros.GET("/:id", buscarLivroPorID)
		livros.POST("", criarLivro)
		livros.PUT("/:id", atualizarLivro)
		livros.DELETE("/:id", deletarLivro)
	}
}

// listarLivros lista todos os livros.
func listarLivros(c *gin.Context) {
	ctx := c.Request.Context()

	livros, err := service.ListarLivros(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar livros"})
		return
	}

	c.JSON(http.StatusOK, livros)
}

// buscarLivroPorID busca um livro pelo seu ID.
func buscarLivroPorID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	livro, err := service.BuscarLivroPorID(ctx, uint(id))
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

// criarLivro cria um novo livro.
func criarLivro(c *gin.Context) {
	ctx := c.Request.Context()

	var novoLivro models.Livro
	if err := c.ShouldBindJSON(&novoLivro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos"})
		return
	}

	if err := service.CriarLivro(ctx, &novoLivro); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar livro"})
		return
	}

	c.JSON(http.StatusCreated, novoLivro)
}

// atualizarLivro atualiza um livro existente.
func atualizarLivro(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	var livroAtualizado models.Livro
	if err := c.ShouldBindJSON(&livroAtualizado); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos"})
		return
	}

	livro, err := service.AtualizarLivro(ctx, uint(id), &livroAtualizado)
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

// deletarLivro remove um livro.
func deletarLivro(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	if err := service.DeletarLivro(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao deletar livro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Livro deletado com sucesso"})
}
