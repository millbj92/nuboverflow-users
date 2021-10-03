package main

import (
	"log"

	"github.com/millbj92/nuboverflow-users/internal/repository"
	"github.com/millbj92/nuboverflow-users/internal/transport/http"
	user "github.com/millbj92/nuboverflow-users/internal/user/service"
)

func Run() error {
	userStore, err := repository.New()
	if err != nil {
		return err
	}

	userService := user.NewService(userStore)
	app := http.CreateRoutes(userService)
	if err != nil {
		return err
	}

	if err := app.Listen(":3000"); err != nil {
		return err
	}
	log.Print("App listening on port 3000.")
	return nil
}



func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
