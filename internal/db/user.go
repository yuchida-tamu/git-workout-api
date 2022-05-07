package db

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	appUser "github.com/yuchida-tamu/git-workout-api/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
}

type UserRow struct {
	ID       string
	Username string
	Password string
}

func convertUserRowToUser(row UserRow) appUser.User {
	return appUser.User{
		ID:       row.ID,
		Username: row.Username,
		Password: row.Password,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (d *Database) GetUsers(ctx context.Context) ([]appUser.User, error) {
	var users []appUser.User
	rows, err := d.Client.QueryContext(
		ctx,
		`SELECT * 
		FROM users`,
	)

	if err != nil {
		return []appUser.User{}, fmt.Errorf("error fetching the user: %w", err)
	}

	for rows.Next() {
		var ID, Username, Password string

		err := rows.Scan(&ID, &Username, &Password)
		if err != nil {
			return []appUser.User{}, fmt.Errorf("error fetching the user: %w", err)
		}

		users = append(users, convertUserRowToUser(UserRow{ID: ID, Username: Username, Password: Password}))

	}

	return users, nil
}

func (d *Database) GetUser(ctx context.Context, uuid string) (appUser.User, error) {
	var userRow UserRow

	row := d.Client.QueryRowContext(
		ctx,
		`SELECT * FROM users WHERE id = $1`,
		uuid,
	)
	err := row.Scan(&userRow.ID, &userRow.Username, &userRow.Password)
	if err != nil {
		return appUser.User{}, fmt.Errorf("error fetching the user by uuid: %w", err)
	}

	return convertUserRowToUser(userRow), nil
}

func (d *Database) GetUserByUsername(ctx context.Context, username string) (appUser.User, error) {
	var userRow UserRow

	row := d.Client.QueryRowContext(
		ctx,
		`SELECT * FROM users WHERE username = $1`,
		username,
	)

	err := row.Scan(&userRow.ID, &userRow.Username, &userRow.Password)
	if err != nil {
		return appUser.User{}, fmt.Errorf("error fetching the user by username: %w", err)
	}

	return convertUserRowToUser(userRow), nil
}

func (d *Database) PostUser(ctx context.Context, user appUser.User) (appUser.User, error) {
	user.ID = uuid.NewV4().String()
	hash, err := hashPassword(user.Password)
	if err != nil {
		return appUser.User{}, err
	}
	user.Password = hash

	postRow := UserRow{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}

	row, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO users
		(id, username, password)
		VALUES
		(:id, :username, :password)`,
		postRow,
	)

	if err != nil {
		return appUser.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	if err := row.Close(); err != nil {
		return appUser.User{}, fmt.Errorf("failed to close row: %w", err)
	}

	return user, nil
}

func (d *Database) UpdateUser(ctx context.Context, uuid string, user appUser.User) (appUser.User, error) {
	userRow := UserRow{
		ID:       uuid,
		Username: user.Username,
		Password: user.Password,
	}

	row, err := d.Client.NamedQueryContext(
		ctx,
		`UPDATE users SET
		username = :username,
		password = : password
		WHERE id = :id`,
		userRow,
	)

	if err != nil {
		return appUser.User{}, fmt.Errorf("failed to update user: %w", err)
	}
	if err := row.Close(); err != nil {
		return appUser.User{}, fmt.Errorf("failed to close row: %w", err)
	}

	return convertUserRowToUser(userRow), nil
}

func (d *Database) DeleteUser(ctx context.Context, uuid string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM users WHERE id = $1`,
		uuid,
	)
	if err != nil {
		return fmt.Errorf("failed to delete user from database: %w", err)
	}

	return nil
}
