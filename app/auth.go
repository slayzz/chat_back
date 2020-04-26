package app

import (
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"todoapp/model"
	u "todoapp/utils"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = u.Message(false, "missing auth token")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "invalid/malformed auth token")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1]
		tk := &model.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv(TokenPassword)), nil
		})

		if err != nil {
			response = u.Message(false, "malformed authentication token")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		if !token.Valid {
			response = u.Message(false, "token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
