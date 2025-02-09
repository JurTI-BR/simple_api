package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Set JWT secret for tests.
	err := os.Setenv("JWT_SECRET", "test_secret")
	if err != nil {
		return
	}

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "NoAuthorizationHeader",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"message":"Token não fornecido"}`,
		},
		{
			name:         "InvalidTokenPrefix",
			authHeader:   "InvalidPrefix tokenstring",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"message":"Token inválido"}`,
		},
		{
			name:         "EmptyToken",
			authHeader:   "Bearer ",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"message":"Token inválido"}`,
		},
		{
			name:         "InvalidToken",
			authHeader:   "Bearer invalidtoken",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"message":"Token inválido"}`,
		},
		{
			name:         "ValidToken",
			authHeader:   "Bearer " + generateValidToken(),
			expectedCode: http.StatusOK,
			expectedBody: "OK",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a Gin engine.
			r := gin.New()

			// Add the AuthMiddleware.
			r.Use(AuthMiddleware())
			// Add a dummy handler to simulate the protected route.
			r.GET("/protected", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			// Create a request.
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", test.authHeader)

			// Record the response.
			w := httptest.NewRecorder()

			// Process the request using the Gin engine.
			r.ServeHTTP(w, req)

			// Assertions.
			assert.Equal(t, test.expectedCode, w.Code)

			// Compare response body only for unauthorized cases.
			if w.Code != http.StatusOK {
				assert.JSONEq(t, test.expectedBody, w.Body.String())
			} else {
				assert.Equal(t, test.expectedBody, w.Body.String())
			}
		})
	}
}

// generateValidToken creates a valid JWT for testing purposes.
func generateValidToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test_user",
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}
