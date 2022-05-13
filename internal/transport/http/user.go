package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yuchida-tamu/git-workout-api/internal/user"
)

type PostUserRequest struct {
	Username string
	Password string
}

type AuthData struct {
	Username string
	Password string
}

type UserForClient struct {
	ID       string
	Username string
}

type AuthUserResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id`
}

type UserService interface {
	GetUser(ctx context.Context, ID string) (user.User, error)
	PostUser(context.Context, user.User) (user.User, error)
	UpdateUser(ctx context.Context, ID string, user user.User) (user.User, error)
	DeleteUser(ctx context.Context, ID string) error
	AuthUser(ctx context.Context, username string, password string) (user.User, error)
}

// TODO: remove password from reponse

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

	userForClient := UserForClient{
		ID:       postedUser.ID,
		Username: postedUser.Username,
	}

	if err := json.NewEncoder(w).Encode(userForClient); err != nil {
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

	userForClient := UserForClient{
		ID:       user.ID,
		Username: user.Username,
	}

	if err := json.NewEncoder(w).Encode(userForClient); err != nil {
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

	userForClient := UserForClient{
		ID:       user.ID,
		Username: user.Username,
	}

	if err := json.NewEncoder(w).Encode(userForClient); err != nil {
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

func (h *Handler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var authData AuthData
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(authData)

	if err != nil {
		http.Error(w, "not a valid input", http.StatusBadRequest)
		return
	}

	user, err := h.Service.User.AuthUser(r.Context(), authData.Username, authData.Password)

	if err != nil {
		response := AuthUserResponse{
			Token:  "",
			UserID: "",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	// create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID          works too)
	t, err := token.SignedString([]byte(os.Getenv("SIGNING_SECRET")))
	if err != nil {
		response := AuthUserResponse{
			Token:  "",
			UserID: "",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	response := AuthUserResponse{
		Token:  t,
		UserID: user.ID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}

}
