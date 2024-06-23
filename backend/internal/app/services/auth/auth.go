package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"kinogo/pkg/db"
	"log"
	"sort"
	"strings"
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) ValidateTelegramAuth(data map[string]string, botToken string) bool {
	checkHash := data["hash"]
	delete(data, "hash")

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var checkString strings.Builder
	for _, key := range keys {
		if value := data[key]; value != "" {
			checkString.WriteString(key + "=" + value + "\n")
		}
	}
	checkStringString := strings.TrimSuffix(checkString.String(), "\n")

	// Создаем SHA256 хеш от токена бота
	botTokenHash := sha256.Sum256([]byte(botToken))
	secretKey := botTokenHash[:]

	hmacHash := hmac.New(sha256.New, secretKey)
	hmacHash.Write([]byte(checkStringString))
	calculatedHash := hex.EncodeToString(hmacHash.Sum(nil))

	fmt.Println(calculatedHash)

	return strings.EqualFold(calculatedHash, checkHash)
}

func (s Service) AddUserIfNotExists(data map[string]string) {
	id := data["id"]

	// Проверка на существование пользователя
	var exists bool
	err := db.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		fmt.Errorf("error checking user existence: %v", err)
	}

	if exists {
		// Пользователь уже существует
		log.Printf("User with ID %s already exists.", id)
	}

	username := data["username"]
	photoURL := data["photo_url"]
	firstName := data["first_name"]
	lastName := data["last_name"]
	authDate := data["auth_date"]

	// Добавление пользователя, если его нет
	_, err = db.Conn.Exec(`INSERT INTO users (id, username, photourl, first_name, last_name, auth_date) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, username, photoURL, firstName, lastName, authDate)
	if err != nil {
		fmt.Errorf("error inserting new user: %v", err)
	}

	log.Printf("User with ID %s added successfully.", id)
	return
}

func (s Service) GenerateToken(userID string, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}
