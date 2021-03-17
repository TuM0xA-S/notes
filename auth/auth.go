package auth

import (
	"context"
	"net/http"
	. "notes/config"
	"notes/util"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

//UserID allows to get id of authorized user
var UserID = &struct{}{}

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
			response := util.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			util.Respond(w, response)
			return
		}

		splitted := strings.Fields(tokenHeader)
		if len(splitted) != 2 {
			response := util.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			util.Respond(w, response)
			return
		}

		tokenPart := splitted[1]

		tk := &Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(Cfg.TokenPassword), nil
			})

		if err != nil {
			response := util.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			util.Respond(w, response)
			return
		}

		if !token.Valid {
			response := util.Message(false, "Token is not valid")
			w.WriteHeader(http.StatusForbidden)
			util.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), UserID, tk.UserID)
		r = r.WithContext(ctx)
		hand.ServeHTTP(w, r)
	})
}
