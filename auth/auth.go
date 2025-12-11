package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Key = []byte("eifbwrrpbep")

func CreateToken(user int) string {
	claim := jwt.MapClaims{
		"user": user,
		"ext":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	finalToken, _ := token.SignedString(Key)
	return finalToken
}

func Check(w http.ResponseWriter, r *http.Request) (int, bool) {
	token := r.Header.Get("Authorization")
	finalToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if err != nil || !finalToken.Valid {
		fmt.Fprintln(w, "bad token")
		return 0, false
	}
	claims := finalToken.Claims.(jwt.MapClaims)
	userID := int(claims["user"].(float64))
	return userID, true
}
