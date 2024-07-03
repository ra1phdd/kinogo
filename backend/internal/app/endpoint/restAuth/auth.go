package restAuth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kinogo/pkg/logger"
	"net/http"
)

type Auth interface {
	ValidateTelegramAuth(data map[string]interface{}, botToken string) bool
	AddUserIfNotExists(data map[string]interface{})
	GenerateToken(data map[string]interface{}, jwtSecret string) (string, error)
	ValidateToken(tokenString string, data map[string]interface{}, jwtSecret string) (bool, error)
	CheckAdminService(id int32) (bool, error)
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

		isAdmin, err := e.Auth.CheckAdminService(int32(authData.ID))
		if err != nil {
			logger.Error("Ошибка в работе функции CheckAuthService", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}

		// Преобразуем authData в map для проверки подписи
		data := map[string]interface{}{
			"id":         authData.ID,
			"first_name": authData.FirstName,
			"last_name":  authData.LastName,
			"username":   authData.Username,
			"photo_url":  authData.PhotoURL,
			"auth_date":  authData.AuthDate,
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

		delete(data, "hash")
		data["isAdmin"] = isAdmin

		token, err := e.Auth.GenerateToken(data, jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"token":     token,
			"userId":    fmt.Sprint(authData.ID),
			"firstName": authData.FirstName,
			"lastName":  authData.LastName,
			"username":  authData.Username,
			"photoUrl":  authData.PhotoURL,
			"authDate":  fmt.Sprint(authData.AuthDate),
			"isAdmin":   isAdmin,
		})
	}
}
