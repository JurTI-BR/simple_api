package service

import (
	"books_api/models"
	"books_api/repository"
	"context"
)

func ListarLivros(ctx context.Context) ([]models.Livro, error) {
	return repository.GetLivrosFromCache(ctx)
}
