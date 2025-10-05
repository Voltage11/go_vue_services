package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Структура пользовательских данных
type UserClaims struct {
	ID    uint    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

// Создание JWT токена из структуры
func CreateToken(user UserClaims, jwtKey []byte) (string, error) {
	// Устанавливаем время жизни токена
	user.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 часа
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "services",
		Subject:   "user-auth",
	}

	// Создаем токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)

	// Подписываем токен ключом
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Валидация и получение структуры из JWT токена
func ParseToken(tokenString string, jwtKey []byte) (*UserClaims, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Извлекаем claims
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}