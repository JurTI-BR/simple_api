package service

import (
	"context"
	"fmt"

	"books_api/models"
	"books_api/repository"
)

// ListarLivros retorna todos os livros, utilizando o cache quando disponível.
// Caso ocorra algum erro, ele é propagado com contexto adicional.
func ListarLivros(ctx context.Context) ([]models.Livro, error) {
	livros, err := repository.GetLivrosFromCache(ctx, 1, 100, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar livros: %w", err)
	}
	return livros, nil
}

// BuscarLivroPorID retorna um livro com base no ID fornecido.
// Caso não seja encontrado ou ocorra algum erro, ele é propagado com contexto adicional.
func BuscarLivroPorID(ctx context.Context, id uint) (*models.Livro, error) {
	livro, err := repository.GetLivroByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar livro com ID %d: %w", id, err)
	}
	return livro, nil
}

// CriarLivro insere um novo livro no banco de dados.
// Em caso de erro na inserção, o erro é propagado com contexto adicional.
func CriarLivro(ctx context.Context, livro *models.Livro) error {
	if err := repository.CreateLivro(ctx, livro); err != nil {
		return fmt.Errorf("erro ao criar livro: %w", err)
	}
	return nil
}

// AtualizarLivro atualiza os dados de um livro existente com base no ID.
// O livro atualizado é retornado ou, em caso de erro, ele é propagado com contexto adicional.
func AtualizarLivro(ctx context.Context, id uint, livroAtualizado *models.Livro) (*models.Livro, error) {
	livro, err := repository.UpdateLivro(ctx, id, livroAtualizado)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar livro com ID %d: %w", id, err)
	}
	return livro, nil
}

// DeletarLivro remove um livro do banco de dados com base no ID fornecido.
// Se ocorrer algum erro durante a remoção, ele é propagado com contexto adicional.
func DeletarLivro(ctx context.Context, id uint) error {
	if err := repository.DeleteLivro(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar livro com ID %d: %w", id, err)
	}
	return nil
}
