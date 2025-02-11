package routes

import (
	"books_api/middleware"
	"books_api/models"
	"books_api/service"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BookRoutes(router *gin.Engine, livroService service.LivroService) { // Receber o serviço corretamente
	livros := router.Group("/livros")
	livros.Use(middleware.AuthMiddleware())
	{
		livros.GET("", func(c *gin.Context) { listarLivros(c, livroService) })
		livros.GET("/:id", func(c *gin.Context) { buscarLivroPorID(c, livroService) })
		livros.POST("", func(c *gin.Context) { criarLivro(c, livroService) })
		livros.PUT("/:id", func(c *gin.Context) { atualizarLivro(c, livroService) })
		livros.DELETE("/:id", func(c *gin.Context) { deletarLivro(c, livroService) })
		livros.POST("/:id/upload", func(c *gin.Context) { uploadImagemLivro(c, livroService) }) // Passando o serviço
	}
}

func uploadImagemLivro(c *gin.Context, srv service.LivroService) { // Agora recebe o serviço
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Arquivo inválido"})
		return
	}

	imageDir := "uploads/"
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		os.Mkdir(imageDir, os.ModePerm)
	}

	filename := fmt.Sprintf("%d%s", id, filepath.Ext(file.Filename))
	imagePath := filepath.Join(imageDir, filename)
	fullImagePath := "uploads/" + filename // Caminho salvo no banco

	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao salvar imagem"})
		return
	}

	// Atualizar a imagem no banco de dados usando o serviço corretamente
	if err := srv.AtualizarImagemLivro(id, fullImagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar imagem no banco"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagem enviada com sucesso!", "imagePath": fullImagePath})
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

	// Adicionar URL completa para a imagem
	for i := range livros {
		if livros[i].ImagePath != "" {
			livros[i].ImagePath = "http://localhost:8080/" + livros[i].ImagePath
		}
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
