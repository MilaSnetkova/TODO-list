package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/MilaSnetkova/TODO-list/internal/config"
)

// секретный ключ для подписания
var jwtSecret = []byte("your-secret-key")  

func SignInHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error":"Failed to decode JSON"}`, http.StatusBadRequest)
		return
	}

	currentPassword := cfg.Password
	if request.Password != currentPassword {
		// Если пароль неверен, возвращаем ошибку
		http.Error(w, `{"error": "Wrong password"}`, http.StatusUnauthorized)
		return
	}


	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Время жизни токена 8 часов
		"exp": time.Now().Add(8 * time.Hour).Unix(),  
	})

	// Подписываем токен
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, `{"error":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"token": tokenString})
}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Проверяем токен
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}