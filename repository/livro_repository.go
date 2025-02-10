package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
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

// GetLivrosFromCache retorna uma lista paginada de livros, tentando primeiro obter os dados do cache.
// Caso não haja cache ou ocorra erro, os dados são buscados no banco de dados e o cache é atualizado.
func GetLivrosFromCache(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.Livro, error) {
	var livros []models.Livro
	offset := (page - 1) * limit

	cacheKeyWithParams := cacheKey + getCacheKey(filters, page, limit)

	// Tenta obter os livros do cache Redis.
	data, err := config.RedisClient.Get(ctx, cacheKeyWithParams).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(data), &livros); err == nil {
			return livros, nil
		}
		log.Printf("Erro ao desserializar livros do cache: %v", err)
	} else if err != redis.Nil {
		log.Printf("Erro ao buscar livros no Redis: %v", err)
	}

	// Monta a query com os filtros
	query := config.DB.WithContext(ctx).Model(&models.Livro{})
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Busca os livros no banco de dados
	if err = query.Offset(offset).Limit(limit).Find(&livros).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Nenhum livro encontrado no banco de dados")
			return nil, nil
		}
		log.Printf("Erro ao buscar livros no banco de dados: %v", err)
		return nil, err
	}

	// Atualiza o cache Redis com os dados obtidos.
	if err := updateCache(ctx, cacheKeyWithParams, livros); err != nil {
		log.Printf("Erro ao atualizar cache: %v", err)
	}

	return livros, nil
}

// GetLivroByID retorna um livro pelo seu ID (com cache).
func GetLivroByID(ctx context.Context, id uint) (*models.Livro, error) {
	cacheKeyByID := cacheKey + ":id:" + strconv.Itoa(int(id))

	var livro models.Livro

	// Tenta buscar do cache
	data, err := config.RedisClient.Get(ctx, cacheKeyByID).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(data), &livro); err == nil {
			return &livro, nil
		}
		log.Printf("Erro ao desserializar livro do cache: %v", err)
	} else if err != redis.Nil {
		log.Printf("Erro ao acessar Redis: %v", err)
	}

	// Busca no banco de dados
	if err := config.DB.WithContext(ctx).First(&livro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Atualiza cache
	go func() {
		if err := updateCache(ctx, cacheKeyByID, livro); err != nil {
			log.Printf("Erro ao atualizar cache de livro por ID: %v", err)
		}
	}()

	return &livro, nil
}

// CreateLivro adiciona um novo livro ao banco de dados dentro de uma transação e invalida o cache.
func CreateLivro(ctx context.Context, livro *models.Livro) error {
	return config.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(livro).Error; err != nil {
			return err
		}

		// Invalida cache após inserção
		go func() {
			err := invalidateCache(ctx)
			if err != nil {

			}
		}()
		return nil
	})
}

// UpdateLivro atualiza um livro existente dentro de uma transação e invalida o cache.
func UpdateLivro(ctx context.Context, id uint, livroAtualizado *models.Livro) (*models.Livro, error) {
	var livro models.Livro

	err := config.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&livro, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		livro.Titulo = livroAtualizado.Titulo
		livro.Autor = livroAtualizado.Autor
		if livroAtualizado.ImagePath != "" {
			livro.ImagePath = livroAtualizado.ImagePath
		}

		if err := tx.Save(&livro).Error; err != nil {
			return err
		}

		// Invalida cache em segundo plano
		go func() {
			err := invalidateCache(ctx)
			if err != nil {
				log.Printf("Erro ao invalidar cache: %v", err)
			}
		}()
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &livro, nil
}

// DeleteLivro remove um livro dentro de uma transação e invalida o cache.
func DeleteLivro(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Livro{}, id).Error; err != nil {
			return err
		}

		go func() {
			err := invalidateCache(ctx)
			if err != nil {

			}
		}()
		return nil
	})
}

// updateCache armazena os livros no Redis com um tempo de expiração definido.
func updateCache(ctx context.Context, key string, data interface{}) error {
	cacheData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, key, cacheData, cacheExpiration).Err()
}

// invalidateCache remove a chave do cache para garantir dados atualizados.
func invalidateCache(ctx context.Context) error {
	return config.RedisClient.Del(ctx, cacheKey).Err()
}

// getCacheKey gera uma chave de cache personalizada para filtros e paginação.
func getCacheKey(filters map[string]interface{}, page, limit int) string {
	key := ""
	for k, v := range filters {
		key += ":" + k + "=" + v.(string)
	}
	return key + ":page=" + string(rune(page)) + ":limit=" + string(rune(limit))
}
