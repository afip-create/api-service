package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	AuthTokenPrefix = "Bearer "
)

type JWTPayload struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email   string `json:"email"`
	Role    string `json:"role"`
}

func GenerateToken(username, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTPayload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
		Username: username,
		Email:   email,
		Role:    role,
	})
	return token.SignedString([]byte("secret"))
}

func ValidateToken(token string) (string, error) {
	token = strings.TrimPrefix(token, AuthTokenPrefix)
	_, err := jwt.ParseWithClaims(token, &JWTPayload{}, func(token *jwt.Token) (interface{}, error {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	return token.Claims.(*JWTPayload).Username, nil
}

func GetUsernameFromToken(token string) (string, error) {
	username, err := ValidateToken(token)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetEmailFromToken(token string) (string, error) {
	var payload JWTPayload
	_, err := jwt.ParseWithClaims(token, &payload, func(token *jwt.Token) (interface{}, error {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	return payload.Email, nil
}

func GetRoleFromToken(token string) (string, error) {
	var payload JWTPayload
	_, err := jwt.ParseWithClaims(token, &payload, func(token *jwt.Token) (interface{}, error {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	return payload.Role, nil
}

func ValidateAdminToken(token string) (bool, error) {
	role, err := GetRoleFromToken(token)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, AuthTokenPrefix) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		username, err := GetUsernameFromToken(token)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Request took %v", time.Since(start))
	})
}

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}