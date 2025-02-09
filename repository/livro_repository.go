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

// GetLivrosFromCache tenta recuperar os livros do cache Redis.
// Caso não encontre ou ocorra algum erro, busca os dados no banco de dados,
// atualiza o cache e retorna os livros.
func GetLivrosFromCache(ctx context.Context) ([]models.Livro, error) {
	const cacheKey = "livros"
	var livros []models.Livro

	// Tenta obter os livros do cache Redis.
	data, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache encontrado: realiza a desserialização dos dados.
		if err = json.Unmarshal([]byte(data), &livros); err == nil {
			return livros, nil
		}
		log.Printf("Erro ao desserializar livros do cache: %v", err)
	} else if err != redis.Nil {
		// Erro inesperado ao acessar o Redis (diferente de cache miss).
		log.Printf("Erro ao buscar livros no Redis: %v", err)
	}

	// Caso não haja cache (ou erro na desserialização), busca os livros no banco de dados.
	if err = config.DB.Find(&livros).Error; err != nil {
		// Caso não haja registros, pode ser retornado nil.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Nenhum livro encontrado no banco de dados")
			return nil, nil
		}
		log.Printf("Erro ao buscar livros no banco de dados: %v", err)
		return nil, err
	}

	// Serializa os livros para armazená-los no cache.
	cacheData, err := json.Marshal(livros)
	if err != nil {
		log.Printf("Erro ao serializar livros para cache: %v", err)
		return livros, nil // Retorna os dados mesmo que o cache não seja atualizado.
	}

	// Atualiza o cache Redis com um tempo de expiração de 10 minutos.
	if err = config.RedisClient.Set(ctx, cacheKey, cacheData, 10*time.Minute).Err(); err != nil {
		log.Printf("Erro ao armazenar livros no cache Redis: %v", err)
	}

	return livros, nil
}
