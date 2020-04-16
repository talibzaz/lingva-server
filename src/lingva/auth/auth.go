package auth

import (
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"encoding/json"
	"github.com/gorilla/context"
	"time"
	"fmt"
)

func GenerateToken() (string, error){
	authKey := "LingvaServerAuthKey"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "username",
		"date": time.Date(2018, 9, 1, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString([]byte(authKey))
	if err != nil {
		fmt.Println("error in token", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		authKey := "LingvaServerAuthKey"

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token)(interface{}, error){
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, errors.New("there was an error")
					}
					return []byte(authKey), nil
				})
				if err != nil {
					json.NewEncoder(w).Encode(err.Error())
					return
				}
				if token.Valid {
					context.Set(r, "decoded", token.Claims)
					next(w, r)
				} else {
					json.NewEncoder(w).Encode("Invalid authorizaztion token")
				}
			}
		} else {
			json.NewEncoder(w).Encode("An authorization header is required")
		}
	})
}
