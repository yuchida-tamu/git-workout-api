package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password string
}

type Store interface {
	GetUsers(context.Context) ([]User, error)
	GetUser(context.Context, string) (User, error)
	PostUser(context.Context, User) (User, error)
	UpdateUser(context.Context, string, User) (User, error)
	DeleteUser(context.Context, string) error
	GetUserByUsername(context.Context, string) (User, error)
}

type Service struct {
	Store Store
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.Store.GetUsers(ctx)
	if err != nil {
		fmt.Println(err)
		return []User{}, err
	}
	return users, nil
}

func (s *Service) GetUser(ctx context.Context, ID string) (User, error) {
	user, err := s.Store.GetUser(ctx, ID)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return user, nil
}

func (s *Service) PostUser(ctx context.Context, user User) (User, error) {
	user, err := s.Store.PostUser(ctx, user)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, ID string, user User) (User, error) {
	user, err := s.Store.UpdateUser(ctx, ID, user)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, ID string) error {
	err := s.Store.DeleteUser(ctx, ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Service) AuthUser(ctx context.Context, username string, password string) (User, error) {
	user, err := s.Store.GetUserByUsername(ctx, username)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	if !checkPasswordHash(password, user.Password) {
		return User{}, fmt.Errorf("Failed to authenticate the user")
	}

	return user, nil
}
