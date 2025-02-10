package config

import (
	"fmt"
	"log"
	"os"

	"books_api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase inicializa a conexão com o banco de dados e executa as migrações necessárias.
func ConnectDatabase() {
	// Obtém as variáveis de ambiente necessárias para a conexão.
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Verifica se as variáveis obrigatórias estão configuradas.
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("Variáveis de ambiente do banco de dados não configuradas corretamente")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Executa as migrações dos modelos.
	if err = DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Erro ao migrar o modelo User: %v", err)
	}
	if err = DB.AutoMigrate(&models.Livro{}); err != nil {
		log.Fatalf("Erro ao migrar o modelo Livro: %v", err)
	}

	log.Println("Banco de dados conectado e tabelas migradas com sucesso!")
}
