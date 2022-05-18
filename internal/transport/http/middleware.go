package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// JSONMiddleware - a middleware function to set Http Header
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			},
		).Info("handled request")

		next.ServeHTTP(w, r)
	})

}

// TimeoutMiddleware - middleware function to timeout after 15 seconds.
func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AddCurrentUserToContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve user Id from jwt
		t := RetrieveJWTTokenFromHeader(r)
		if t == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SIGNING_SECRET")), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userIdInToken := tokenClaims["userId"]
		if userIdInToken == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// add user id to the current context
		ctx := context.WithValue(r.Context(), "userId", userIdInToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
