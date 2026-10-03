package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token,omitempty"`
}

// passwordHash возвращает SHA-256 хэш пароля в виде hex-строки
func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// signInHandler обрабатывает POST /api/signin — проверку пароля и выдачу JWT
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "ошибка десериализации JSON")
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeError(w, "Неверный пароль")
		return
	}

	claims := jwt.MapClaims{
		"hash": passwordHash(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(pass))
	if err != nil {
		writeError(w, "ошибка создания токена")
		return
	}

	writeJSON(w, signInResponse{Token: signedToken})
}

// auth — middleware, проверяющее JWT-токен из куки, если установлен TODO_PASSWORD
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtStr string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookie.Value
			}

			valid := false

			token, err := jwt.Parse(jwtStr, func(t *jwt.Token) (any, error) {
				return []byte(pass), nil
			})

			if err == nil && token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok {
					hash, ok := claims["hash"].(string)
					if ok && hash == passwordHash(pass) {
						valid = true
					}
				}
			}

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	}
}
