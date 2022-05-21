package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		// Bearer: token-string, (in the Header, Request should have "Authorization", which is formatted as "bearer [encoded jwt key]")
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		valid, expired := validateToken(authHeaderParts[1])
		if !valid {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		if expired {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JWTAuth(
	original func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		// Bearer: token-string, (in the Header, Request should have "Authorization", which is formatted as "bearer [encoded jwt key]")
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		valid, expired := validateToken(authHeaderParts[1])
		if expired {
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(ErrorResponse{
				Message:   "token expired",
				ErrorCode: -2000,
			})
			if err != nil {
				http.Error(w, "something went wrong", http.StatusUnauthorized)
				return
			}
		}
		if !valid {
			http.Error(w, "not authorized invalid token", http.StatusUnauthorized)
			return
		}

		original(w, r)

	}
}

func validateToken(accessToken string) (valid bool, expired bool) {
	var mySigningKey = []byte(os.Getenv("SIGNING_SECRET"))
	t, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}

		return mySigningKey, nil
	})

	v, ok := err.(*jwt.ValidationError)

	if !ok {
		return t.Valid, false
	}

	if v.Errors == jwt.ValidationErrorExpired {
		return false, true
	}

	if err != nil {
		return false, false
	}

	return t.Valid, false
}
