package http

import (
	"net/http"
	"strings"
)

func RetrieveJWTTokenFromHeader(r *http.Request) string {
	authHeader := r.Header["Authorization"]
	if authHeader == nil {
		return ""
	}

	// Bearer: token-string, (in the Header, Request should have "Authorization", which is formatted as "bearer [encoded jwt key]")
	authHeaderParts := strings.Split(authHeader[0], " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return ""
	}

	return authHeaderParts[1]
}
