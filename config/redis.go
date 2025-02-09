package config

import (
	_ "context"
	_ "log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis inicializa a conexão com o Redis
func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // Defina aqui a senha caso necessário
		DB:       0,  // Banco de dados padrão do Redis
	})

}
