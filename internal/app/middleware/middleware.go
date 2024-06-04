package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kinogo/pkg/auth"
	"kinogo/pkg/logger"
	"net/http"
)

func TheLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func CSPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "frame-ancestors http://127.0.0.1")
		next.ServeHTTP(w, r)
	})
}

func CORSPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		verify := auth.VerifyTelegramData(&auth.Auth, auth.BotToken)
		if !verify {
			logger.Warn("Зафиксирована попытка поддельной авторизации")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		logger.Info("Пользователь авторизован")
		c.Next()
	}
}

func NotFoundCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validPaths := map[string]bool{
			"/":          true,
			"/filter":    true,
			"/search":    true,
			"/films":     true,
			"/cartoons":  true,
			"/telecasts": true,
		}

		if _, exists := validPaths[r.URL.Path]; !exists {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
