package service

import (
	"books_api/models"
	"books_api/repository"
	"context"
	"fmt"
)

// Interface para facilitar o mock nos testes
type LivroService interface {
	ListarLivros(ctx context.Context) ([]models.Livro, error)
	BuscarLivroPorID(ctx context.Context, id uint) (*models.Livro, error)
	CriarLivro(ctx context.Context, livro *models.Livro) error
	AtualizarLivro(ctx context.Context, id uint, livroAtualizado *models.Livro) (*models.Livro, error)
	DeletarLivro(ctx context.Context, id uint) error
}

type livroService struct{}

func NewLivroService() LivroService {
	return &livroService{}
}

// Implementação real do serviço
func (s *livroService) ListarLivros(ctx context.Context) ([]models.Livro, error) {
	livros, err := repository.GetLivrosFromCache(ctx, 1, 100, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar livros: %w", err)
	}
	return livros, nil
}

func (s *livroService) BuscarLivroPorID(ctx context.Context, id uint) (*models.Livro, error) {
	livro, err := repository.GetLivroByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar livro com ID %d: %w", id, err)
	}
	return livro, nil
}

func (s *livroService) CriarLivro(ctx context.Context, livro *models.Livro) error {
	if err := repository.CreateLivro(ctx, livro); err != nil {
		return fmt.Errorf("erro ao criar livro: %w", err)
	}
	return nil
}

func (s *livroService) AtualizarLivro(ctx context.Context, id uint, livroAtualizado *models.Livro) (*models.Livro, error) {
	livro, err := repository.UpdateLivro(ctx, id, livroAtualizado)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar livro com ID %d: %w", id, err)
	}
	return livro, nil
}

func (s *livroService) DeletarLivro(ctx context.Context, id uint) error {
	if err := repository.DeleteLivro(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar livro com ID %d: %w", id, err)
	}
	return nil
}
