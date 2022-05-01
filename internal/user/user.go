package user

import "context"

type User struct {
	ID       string
	Username string
}

type Store interface {
	GetUsers(context.Context) ([]User, error)
	GetUser(context.Context, string) (User, error)
	PostUser(context.Context, User) (User, error)
	UpdateUser(context.Context, string, User) (User, error)
	DeleteUser(context.Context, string) error
}

type Service struct {
	Store Store
}
