package service

import (
	"books_api/models"
	"books_api/repository"
	"context"
)

// Listar todos os livros com cache
func ListarLivros(ctx context.Context) ([]models.Livro, error) {
	return repository.GetLivrosFromCache(ctx)
}

// Buscar livro por ID
func BuscarLivroPorID(ctx context.Context, id uint) (*models.Livro, error) {
	return repository.GetLivroByID(ctx, id)
}

// Criar um novo livro
func CriarLivro(ctx context.Context, livro *models.Livro) error {
	return repository.CreateLivro(ctx, livro)
}

// Atualizar um livro existente
func AtualizarLivro(ctx context.Context, id uint, livroAtualizado *models.Livro) (*models.Livro, error) {
	return repository.UpdateLivro(ctx, id, livroAtualizado)
}

// Deletar um livro
func DeletarLivro(ctx context.Context, id uint) error {
	return repository.DeleteLivro(ctx, id)
}
