package auth

import (
	"contacts/util"
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

//Token is access token for auth
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// RequireAuth decorator with JWT
func RequireAuth(hand http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			response := util.Message(false, "Missing auth token")
			util.Respond(w, response)
			return
		}

		splitted := strings.Fields(tokenHeader)
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			response := util.Message(false, "Invalid/Malformed auth token")
			util.Respond(w, response)
			return
		}

		tokenPart := splitted[1]

		tk := &Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			response := util.Message(false, "Invalid/Malformed auth token")
			util.Respond(w, response)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			response := util.Message(false, "Token is not valid")
			util.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		hand.ServeHTTP(w, r)
	})
}
