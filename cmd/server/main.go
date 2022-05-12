package main

import (
	"fmt"

	"github.com/yuchida-tamu/git-workout-api/internal/db"
	"github.com/yuchida-tamu/git-workout-api/internal/record"
	transportHttp "github.com/yuchida-tamu/git-workout-api/internal/transport/http"
	"github.com/yuchida-tamu/git-workout-api/internal/user"
)

func Run() error {
	fmt.Println("starting up the application")
	// connect to database
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("failed to connect to the database")
		return err
	}
	// migrate database
	if err := db.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	userService := user.NewService(db)
	recordService := record.NewService(db)
	service := transportHttp.Service{
		User:   userService,
		Record: recordService,
	}

	httpHandler := transportHttp.NewHandler(service)
	if err := httpHandler.Serve(); err != nil {
		return nil
	}

	return nil
}

func main() {
	fmt.Println("Hello Wolrd")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
