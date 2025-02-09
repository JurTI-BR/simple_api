package service

import (
	"books_api/models"
	"books_api/repository"
)

func ListarLivros() ([]models.Livro, error) {
	return repository.GetLivrosFromCache()
}
