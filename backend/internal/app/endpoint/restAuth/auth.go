package restAuth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth interface {
	ValidateTelegramAuth(data map[string]string, botToken string) bool
	AddUserIfNotExists(data map[string]string)
	GenerateToken(userID string, jwtSecret string) (string, error)
}

type Endpoint struct {
	Auth Auth
}

type TelegramAuthData struct {
	ID        int    `json:"id" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int    `json:"auth_date" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}

func (e Endpoint) TelegramAuthCallback(jwtSecret string, botToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authData TelegramAuthData
		if err := c.ShouldBindJSON(&authData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		// Преобразуем authData в map для проверки подписи
		data := map[string]string{
			"id":         fmt.Sprint(authData.ID),
			"first_name": authData.FirstName,
			"last_name":  authData.LastName,
			"username":   authData.Username,
			"photo_url":  authData.PhotoURL,
			"auth_date":  fmt.Sprint(authData.AuthDate),
			"hash":       authData.Hash,
		}

		if !e.Auth.ValidateTelegramAuth(data, botToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid Telegram Auth",
			})
			return
		}

		e.Auth.AddUserIfNotExists(data)

		token, err := e.Auth.GenerateToken(fmt.Sprint(authData.ID), jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"token":   token,
		})
	}
}
