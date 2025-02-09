package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"books_api/config"
	"books_api/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	cacheKey        = "livros"
	cacheExpiration = 10 * time.Minute
)

// GetLivrosFromCache retorna a lista de livros, tentando primeiro obter os dados do cache.
// Caso não haja cache ou ocorra algum erro, os dados são buscados no banco de dados e o cache é atualizado.
func GetLivrosFromCache(ctx context.Context) ([]models.Livro, error) {
	var livros []models.Livro

	// Tenta obter os livros do cache Redis.
	data, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Desserializa os dados do cache.
		if err = json.Unmarshal([]byte(data), &livros); err == nil {
			return livros, nil
		}
		log.Printf("Erro ao desserializar livros do cache: %v", err)
	} else if err != redis.Nil {
		// Erro inesperado ao acessar o Redis (diferente de cache miss).
		log.Printf("Erro ao buscar livros no Redis: %v", err)
	}

	// Caso não haja cache ou ocorra erro, busca os livros no banco de dados.
	if err = config.DB.Find(&livros).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Nenhum livro encontrado no banco de dados")
			return nil, nil
		}
		log.Printf("Erro ao buscar livros no banco de dados: %v", err)
		return nil, err
	}

	// Atualiza o cache Redis com os dados obtidos.
	if err := updateCache(ctx, cacheKey, livros); err != nil {
		log.Printf("Erro ao atualizar cache: %v", err)
	}

	return livros, nil
}

// GetLivroByID retorna um livro pelo seu ID.
func GetLivroByID(ctx context.Context, id uint) (*models.Livro, error) {
	var livro models.Livro

	// Busca o livro no banco de dados.
	if err := config.DB.First(&livro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Pode-se optar por retornar um erro específico de "não encontrado".
		}
		return nil, err
	}

	return &livro, nil
}

// CreateLivro adiciona um novo livro ao banco de dados e invalida o cache.
func CreateLivro(ctx context.Context, livro *models.Livro) error {
	if err := config.DB.Create(livro).Error; err != nil {
		return err
	}

	// Invalida o cache após a inserção.
	if err := invalidateCache(ctx); err != nil {
		log.Printf("Erro ao invalidar cache após criar livro: %v", err)
	}

	return nil
}

// UpdateLivro atualiza um livro existente e invalida o cache.
func UpdateLivro(ctx context.Context, id uint, livroAtualizado *models.Livro) (*models.Livro, error) {
	var livro models.Livro

	// Verifica se o livro existe.
	if err := config.DB.First(&livro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Ou retornar um erro específico de "não encontrado".
		}
		return nil, err
	}

	// Atualiza os campos desejados.
	livro.Titulo = livroAtualizado.Titulo
	livro.Autor = livroAtualizado.Autor
	// Caso haja outros campos para atualizar, eles podem ser incluídos aqui.

	if err := config.DB.Save(&livro).Error; err != nil {
		return nil, err
	}

	// Invalida o cache para garantir que dados atualizados sejam retornados nas próximas consultas.
	if err := invalidateCache(ctx); err != nil {
		log.Printf("Erro ao invalidar cache após atualizar livro: %v", err)
	}

	return &livro, nil
}

// DeleteLivro remove um livro do banco de dados e invalida o cache.
func DeleteLivro(ctx context.Context, id uint) error {
	if err := config.DB.Delete(&models.Livro{}, id).Error; err != nil {
		return err
	}

	// Invalida o cache.
	if err := invalidateCache(ctx); err != nil {
		log.Printf("Erro ao invalidar cache após excluir livro: %v", err)
	}

	return nil
}

// updateCache armazena os livros no Redis com um tempo de expiração definido.
func updateCache(ctx context.Context, key string, livros []models.Livro) error {
	cacheData, err := json.Marshal(livros)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, key, cacheData, cacheExpiration).Err()
}

// invalidateCache remove a chave do cache para garantir dados atualizados.
func invalidateCache(ctx context.Context) error {
	return config.RedisClient.Del(ctx, cacheKey).Err()
}
