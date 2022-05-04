package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yuchida-tamu/git-workout-api/internal/user"
)

type UserService interface {
	GetUser(ctx context.Context, ID string) (user.User, error)
	PostUser(context.Context, user.User) (user.User, error)
	UpdateUser(ctx context.Context, ID string, user user.User) (user.User, error)
	DeleteUser(ctx context.Context, ID string) error
}

type Response struct {
	message string
}

type PostUserRequest struct {
	Username string
	Password string
}

func convertPostUserRequestToUser(u PostUserRequest) user.User {
	return user.User{
		Username: u.Username,
		Password: u.Password,
	}
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var user PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		http.Error(w, "not a valid user", http.StatusBadRequest)
		return
	}

	convertedUser := convertPostUserRequestToUser(user)

	postedUser, err := h.Service.User.PostUser(r.Context(), convertedUser)
	if err != nil {
		log.Print(err)
		return
	}

	if err := json.NewEncoder(w).Encode(postedUser); err != nil {
		panic(err)
	}

}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.Service.User.GetUser(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		http.Error(w, "not a valid user", http.StatusBadRequest)
		return
	}

	user, err = h.Service.User.UpdateUser(r.Context(), id, user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.User.DeleteUser(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{message: "Successfully deleted"}); err != nil {
		panic(err)
	}
}
