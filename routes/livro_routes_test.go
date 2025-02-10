package routes

import (
	"books_api/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLivroService struct {
	mockResponse []models.Livro
	mockError    error
}

func (m *mockLivroService) ListarLivros(ctx context.Context) ([]models.Livro, error) {
	return m.mockResponse, m.mockError
}

func TestListarLivros(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		mockResponse []models.Livro
		mockError    error
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "success",
			mockResponse: []models.Livro{{ID: 1, Titulo: "Livro 1"}},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: []models.Livro{{ID: 1, Titulo: "Livro 1"}},
		},
		{
			name:         "internal_error",
			mockResponse: nil,
			mockError:    errors.New("internal"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: gin.H{"message": "Erro ao buscar livros"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criando instância do serviço mockado
			mockService := &mockLivroService{
				mockResponse: tt.mockResponse,
				mockError:    tt.mockError,
			}

			router := gin.Default()
			router.GET("/livros", func(c *gin.Context) {
				livros, err := mockService.ListarLivros(c.Request.Context())
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar livros"})
					return
				}
				c.JSON(http.StatusOK, livros)
			})

			// Criando request fake
			req, _ := http.NewRequest(http.MethodGet, "/livros", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificando status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Verificando body
			var body interface{}
			err := json.Unmarshal(w.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}
