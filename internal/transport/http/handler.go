package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Service struct {
	User   UserService
	Record RecordService
}

type Handler struct {
	Router  *mux.Router
	Service Service
	Server  *http.Server
}

func NewHandler(service Service) *Handler {
	h := &Handler{
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()

	// use middlewares
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	// Auth
	h.Router.HandleFunc("/api/v1/auth/auth", h.AuthUser).Methods("POST")
	h.Router.HandleFunc("/api/v1/auth/refresh", h.RefreshToken).Methods("POST")
	// User
	h.Router.HandleFunc("/api/v1/user", JWTAuth(h.PostUser)).Methods("POST")
	h.Router.HandleFunc("/api/v1/user/{id}", JWTAuth(h.GetUser)).Methods("GET")
	h.Router.HandleFunc("/api/v1/user/{id}", JWTAuth(h.UpdateUser)).Methods("UPDATE")
	h.Router.HandleFunc("/api/v1/user/{id}", JWTAuth(h.DeleteUser)).Methods("DELETE")
	// Record
	h.Router.HandleFunc("/api/v1/record", JWTAuth(h.PostRecord)).Methods("POST")
	h.Router.HandleFunc("/api/v1/record/author/{id}", JWTAuth(h.GetRecordByAuthor)).Methods("GET")
	h.Router.HandleFunc("/api/v1/record/{id}", JWTAuth(h.GetRecordById)).Methods("GET")
	h.Router.HandleFunc("/api/v1/record/{id}", JWTAuth(h.UpdateRecord)).Methods("UPDATE")
	h.Router.HandleFunc("/api/v1/record/{id}", JWTAuth(h.DeleteRecord)).Methods("DELETE")
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	h.Server.Shutdown(ctx)

	return nil
}
