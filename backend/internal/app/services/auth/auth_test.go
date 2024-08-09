package auth_v1

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"kinogo/internal/app/services/testutil"
	"kinogo/pkg/db"
	"testing"
)

func TestService_ValidateTelegramAuth(t *testing.T) {
	s := Service{}

	t.Run("Valid Hash", func(t *testing.T) {
		data := map[string]interface{}{
			"id":        12345,
			"username":  "testuser",
			"auth_date": 1598765432,
			"hash":      "7ba4de2f37608e84700decb6d891d8108d3aa899ee1aa2a339377ec5c9c8a0bf",
		}
		botToken := "testbottoken"

		result := s.ValidateTelegramAuth(data, botToken)
		assert.True(t, result)
	})

	t.Run("Invalid Hash", func(t *testing.T) {
		data := map[string]interface{}{
			"id":        12345,
			"username":  "testuser",
			"auth_date": 1598765432,
			"hash":      "f3a563b22de6400892cde0e74b1e4b1b37c2840b94b22e4e07429c30c7367c2a",
		}
		botToken := "testbottoken"

		result := s.ValidateTelegramAuth(data, botToken)
		assert.False(t, result)
	})
}

func TestService_AddUserIfNotExists(t *testing.T) {
	conn, mockDB, mock, _, _ := testutil.SetupMocks()
	defer mockDB.Close()

	db.Conn = conn

	s := Service{}

	t.Run("User Does Not Exist", func(t *testing.T) {
		data := map[string]interface{}{
			"id":         12345,
			"username":   "testuser",
			"photo_url":  "http://example.com/photo.jpg",
			"first_name": "Test",
			"last_name":  "User",
			"auth_date":  1598765432,
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id=\\$1\\)").WithArgs(data["id"]).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec("INSERT INTO users \\(id, username, photourl, first_name, last_name, auth_date\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)").
			WithArgs(data["id"], data["username"], data["photo_url"], data["first_name"], data["last_name"], data["auth_date"]).WillReturnResult(sqlmock.NewResult(1, 1))

		s.AddUserIfNotExists(data)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		data := map[string]interface{}{
			"id":         12345,
			"username":   "testuser",
			"photo_url":  "http://example.com/photo.jpg",
			"first_name": "Test",
			"last_name":  "User",
			"auth_date":  1598765432,
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id=\\$1\\)").WithArgs(data["id"]).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		s.AddUserIfNotExists(data)
	})
}

func TestService_GenerateToken(t *testing.T) {
	s := Service{}

	t.Run("Generate Token Successfully", func(t *testing.T) {
		data := map[string]interface{}{
			"id":         12345,
			"username":   "testuser",
			"photo_url":  "http://example.com/photo.jpg",
			"first_name": "Test",
			"last_name":  "User",
			"auth_date":  1598765432,
			"isAdmin":    false,
		}
		jwtSecret := "secret"

		token, err := s.GenerateToken(data, jwtSecret)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Parse the token to verify its validity
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})
}
