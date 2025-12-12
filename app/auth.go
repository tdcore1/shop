package app

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *App) CreateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user": userID,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.Key)
}

func (a *App) Check(w http.ResponseWriter, r *http.Request) (int, bool) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "missing token", 401)
		return 0, false
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return a.Key, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "invalid token", 401)
		return 0, false
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user"].(float64))
	return userID, true
}
