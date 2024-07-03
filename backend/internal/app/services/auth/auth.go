package auth_v1

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) ValidateTelegramAuth(data map[string]interface{}, botToken string) bool {
	checkHash, ok := data["hash"].(string)
	if !ok {
		fmt.Println("hash is missing or not a string")
		return false
	}
	delete(data, "hash")

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var checkString strings.Builder
	for _, key := range keys {
		value := fmt.Sprintf("%v", data[key])
		if value == "" {
			continue
		}
		checkString.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	checkStringString := strings.TrimSuffix(checkString.String(), "\n")

	// Создаем SHA256 хеш от токена бота
	botTokenHash := sha256.Sum256([]byte(botToken))
	secretKey := botTokenHash[:]

	hmacHash := hmac.New(sha256.New, secretKey)
	hmacHash.Write([]byte(checkStringString))
	calculatedHash := hex.EncodeToString(hmacHash.Sum(nil))

	if !strings.EqualFold(calculatedHash, checkHash) {
		fmt.Println("calculated hash does not match provided hash")
		return false
	}

	// Проверка времени аутентификации (auth_date)
	authDate, ok := data["auth_date"].(int)
	if !ok {
		fmt.Println("auth_date is missing or not a number")
		return false
	}
	telegramAuthTime := time.Unix(int64(authDate), 0)
	if time.Since(telegramAuthTime) > 24*time.Hour {
		fmt.Println("authentication data is outdated")
		return false
	}

	return true
}

func (s Service) AddUserIfNotExists(data map[string]interface{}) {
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

func (s Service) GenerateToken(data map[string]interface{}, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        data["id"],
		"firstName": data["first_name"],
		"lastName":  data["last_name"],
		"username":  data["username"],
		"photoUrl":  data["photo_url"],
		"authDate":  data["auth_date"],
		"isAdmin":   data["isAdmin"],
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}

func (s Service) ValidateToken(tokenString string, data map[string]interface{}, jwtSecret string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		var id string

		if idFloat, ok := claims["id"].(float64); ok {
			id = strconv.FormatFloat(idFloat, 'f', 0, 64)
		}

		if id != fmt.Sprint(data["id"].(int32)) {
			return false, fmt.Errorf("invalid id claim")
		} else if claims["firstName"] != data["first_name"] {
			return false, fmt.Errorf("invalid firstName claim")
		} else if claims["lastName"] != data["last_name"] {
			return false, fmt.Errorf("invalid lastName claim")
		} else if claims["username"] != data["username"] {
			return false, fmt.Errorf("invalid username claim")
		} else if claims["photoUrl"] != data["photo_url"] {
			return false, fmt.Errorf("invalid photoUrl claim")
		} else if claims["isAdmin"] != data["isAdmin"] {
			return false, fmt.Errorf("invalid isAdmin claim")
		}

		return true, nil
	} else {
		return false, fmt.Errorf("invalid token")
	}
}

func (s Service) CheckAdminService(id int32) (bool, error) {
	authString, err := cache.Rdb.Get(cache.Ctx, fmt.Sprintf("checkAdmin_%d", id)).Result()
	if err == nil && authString != "" {
		switch authString {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("Ошибка при получении данных из Redis")
	}

	rows, err := db.Conn.Query(`SELECT "isAdmin" FROM users WHERE id = $1`, id)
	if err != nil {
		logger.Error("Ошибка выполнения SQL-запроса", zap.Int32("id", id))
		return false, err
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			logger.Error("Ошибка при закрытии rows", zap.Error(errClose))
		}
	}()

	var isAdmin bool
	for rows.Next() {
		errScan := rows.Scan(&isAdmin)
		if errScan != nil {
			logger.Error("Ошибка сканирования строки результата запроса")
			return false, errScan
		}
	}

	err = cache.Rdb.Set(cache.Ctx, fmt.Sprintf("checkAdmin_%d", id), fmt.Sprint(isAdmin), 60*time.Minute).Err()
	if err != nil {
		logger.Error("Ошибка при сохранении данных в Redis")
	}

	return isAdmin, nil
}
